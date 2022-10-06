import { Client, GatewayIntentBits } from 'discord.js';
import { MongoClient } from "mongodb";

import config from './config.json'  assert { type: "json" };

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

const client = new Client({ intents: [GatewayIntentBits.Guilds] });

client.once('ready', () => {
    console.log('Ready!');
});

// IMPORTANT: All replies should be `ephemeral: true`.
client.on('interactionCreate', async interaction => {
    if (!interaction.isChatInputCommand()) return;

    const { commandName } = interaction;

    if (commandName === 'get') {
        const platform = interaction.options.getString('platform');
        if (!platforms.includes(platform.toLowerCase())) {
            return await interaction.reply({ content: `Platform ${platform} is not available.`, ephemeral: true })
        }

        await interaction.reply({ content: "I'm about to send the browser in parts. Download them all.", ephemeral: true });

        const part_paths = await getParts(platform)

        for (let i = 0; i < part_paths.length; i++) {
            await interaction.followUp({ files: [part_paths[i]], ephemeral: true });
        }

        await interaction.followUp({ content: "Open the following file in your browser and drag-n-drop all the parts to bake them into one.", files: ["./onion-rings.html"], ephemeral: true });
    } else if (commandName === 'platforms') {
        await interaction.reply({ content: "Available platforms: " + platforms.map(x => "`" + x + "`").join(', '), ephemeral: true })
    }
});

client.login(config.token);
