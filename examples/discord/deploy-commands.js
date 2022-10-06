// Deletes all slash commands and adds the following:
// /platforms
// /get platform:

import { REST, SlashCommandBuilder, Routes } from 'discord.js';
import config from './config.json'  assert { type: "json" };

const commands = [
	new SlashCommandBuilder()
		.setName('get')
		.setDescription('Sends a Tor Browser build in parts.')
		.addStringOption(option => option.setName('platform').setDescription('The target platform.').setRequired(true)),
	new SlashCommandBuilder()
		.setName('platforms')
		.setDescription('Sends all available platforms.'),
]
	.map(command => command.toJSON());

const rest = new REST({ version: '10' }).setToken(config.token);

rest.put(Routes.applicationCommands(config.clientId), { body: [] })
	.then(() => console.log('Successfully deleted all application commands.'))
	.catch(console.error);

rest.put(Routes.applicationCommands(config.clientId), { body: commands })
	.then((data) => console.log(`Successfully registered ${data.length} application commands.`))
	.catch(console.error);
