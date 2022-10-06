import { MongoClient } from "mongodb";
import slack from '@slack/bolt';
import config from './config.json'  assert { type: "json" };
import { readFileSync } from "fs"
import { basename } from "path"

const { App } = slack

const app = new App({
    signingSecret: config.signingSecret,
    token: config.token,
    socketMode: true,
    appToken: config.appToken
});

const dbClient = new MongoClient(config.mongo);

const database = dbClient.db('kremmidi');
const binaries = database.collection('binaries');
const parts = database.collection('parts');
const platforms = [
    'android-aarch64',
    'android-x86',
    'android-armv7',
    'android-x86_64',
    'linux32',
    'linux64',
    'osx64',
    'win32',
    'win64'
]

/**
 * Returns an array of all parts sorted.
 * @param {string} platform The target platform.
 * @returns {string[]} Array of part paths.
 */
async function getParts(platform) {
    if (!platforms.includes(platform.toLowerCase())) return []

    // Get a binary of platform, sorted by version (biggest first).
    const binary = await binaries.findOne({ platform: platform.toLowerCase() }, { sort: { "version": -1 } });

    if (!binary) return []

    // Get all its parts, sorted by part_no.
    const part_list = await parts.find({ belongs_to: binary._id }).sort({ "part_no": -1 }).toArray()

    // If for some reason the parts returned do not match the
    // amount listed in the binary doc, return.
    if (part_list.length !== binary.parts) return []

    // Return the paths but reverse the order so it's 0 -> -1.
    return part_list.reverse().map(x => x.path)
}

app.message(/\!platforms/i, async ({ say }) => {
    await say("Available platforms: " + platforms.map(x => "`" + x + "`").join(', '));
});

app.message(/\!get (.+)/i, async ({ client, message, say, context }) => {
    const platform = context.matches[1];

    const part_paths = await getParts(platform)

    if (part_paths.length === 0) {
        return await say(`No such platform.`)
    }

    await say(`I'll start sending the binary in *${part_paths.length}* parts. Afterwards I'll send onion-rings with instructions! Please wait for all parts to finish uploading...`)

    for (let i = 0; i < part_paths.length; i++) {
        await client.files.upload({
            channels: message.channel,
            filename: basename(part_paths[i]),
            file: readFileSync(part_paths[i])
        });
    }

    await say(`Please download onion-rings.html and open it in your browser. Afterwards follow the instructions on screen and drag-n-drop all parts into it to get the original binary.`)

    await client.files.upload({
        channels: message.channel,
        filename: "onion-rings.html",
        file: readFileSync("./onion-rings.html")
    });
});

(async () => {
    // Start your app
    await app.start();

    console.log('⚡️ Bolt app is running!');
})();
