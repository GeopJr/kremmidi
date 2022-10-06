import {
    MatrixClient,
    SimpleFsStorageProvider,
    AutojoinRoomsMixin
} from "matrix-bot-sdk";
import { readFileSync } from "fs"
import { basename } from "path"
import { MongoClient } from "mongodb";

import config from './config.json'  assert { type: "json" };

const delay = (timeout) => new Promise((resolve, reject) => setTimeout(resolve, timeout));

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

const homeserverUrl = "https://matrix.org";
const accessToken = config.accessToken;
const storage = new SimpleFsStorageProvider("bot.json");

const client = new MatrixClient(homeserverUrl, accessToken, storage);
AutojoinRoomsMixin.setupOnClient(client);

const botUser = await client.getUserId()

client.on("room.message", async (roomId, event) => {
    if (!event["content"]) return;

    const sender = event["sender"];
    if (sender === botUser) return;

    const body = event["content"]["body"];
    if (!body.startsWith("!")) return;

    const platform = body.replace('!', '').toLowerCase()
    if (!platforms.includes(platform)) return;

    const part_paths = await getParts(platform)
    await client.sendMessage(roomId, {
        "msgtype": "m.notice",
        "body": `I'll start sending the binary in *${part_paths.length}* parts. Afterwards I'll send onion-rings with instructions! Please wait for all parts to finish uploading...`,
    });

    for (let i = 0; i < part_paths.length; i++) {
        const partBuff = readFileSync(part_paths[i])
        const uploadedPart = await client.uploadContent(partBuff);
        await client.sendMessage(roomId, {
            msgtype: "m.file",
            body: basename(part_paths[i]),
            file: {
                url: uploadedPart
            },
        });


        await delay(2500);
    }

    const onionRings = readFileSync("./onion-rings.html")
    const uploadedOnionRings = await client.uploadContent(onionRings);

    await client.sendMessage(roomId, {
        "msgtype": "m.notice",
        "body": "Please download onion-rings.html and open it in your browser. Afterwards follow the instructions on screen and drag-n-drop all parts into it to get the original binary.",
    });

    // Matrix bot api doesn't seem to handle
    // rate limits :/
    await delay(2500);

    await client.sendMessage(roomId, {
        msgtype: "m.file",
        body: "onion-rings.html",
        file: {
            url: uploadedOnionRings
        },
    });
});

client.start().then(() => console.log("Client started!"));