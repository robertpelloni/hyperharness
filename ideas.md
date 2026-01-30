For every link, have a (preferably lightweight and/or free) subagent research it: scrape it, read it, index it, categorize it, summarize it. If it's an exact duplicate skip it. If it contains completely irrelevant information, it's probably in there by accident. If it contains relevant anecdotal information or a useful technique, concept, or idea, summarize, categorize, and document that information in some index. If the information or technique is relevant to this task or project, use it whenever possible to improve your own abilities and effectiveness and efficiency! Add it to AGENTS.md (or other more relevant md file). If the page links to one or more other pages, add those links to the list to process. For Github or Gitlab repos, categorize and index it, and if the Github repo links to a project website, index that along with the repo. Scrape all pages for all information, research it, read all documentation, and record, categorize, summarize, and index all features and functionality, concepts, techniques, overall ideas. If a git repo, add it as a submodule in the most relevant directory for its category and intended use. Every feature is a potential feature to add to aios, either reimplement it and use the submodule simply for reference, or if possible call out to and reuse existing code, or wrap over the code in order to achieve its functionality whenever possible, so that we can immediately and seamlessly benefit from any and all updates to that submodule. 


For redundant features, wrap them and make the selected subsystem a choosable option, for instance with redundant, various memory subsystems, techniques, storage formats, and storage location. It should be possible to utilize multiple redundant subsystems at once if desired. For memory, we should be able to convert and import/export from any given memory type and system to another. Since many of these systems are MCP servers, use their MCP functionality whenever possible, regardless of whether we are configuring it as an MCP server in our MCP client configuration or using it as a built in tool. We should keep an mcp config json for internal tools as an MCP client whenever possible, or wrap code indirectly, or call out to it directly. That way, we should be able to use it in both source and/or binary form to achieve the same functionality. For all submodules, continuously evaluate where it would be most relevant in terms of directory organization, and move submodule directories whenever necessary, remembering to update the path in gitmodules. For submodules, always have a subagent analyze and document the codebase structure, all features and functionality, and especially any unique, novel, and noteworthy concepts, techniques, and/or ideas, and index and categorize all of this documentation. In a sense, our CLI/TUI/WebUI itself is sort of made up of parts or plugins both directly in code form and indirectly as an MCP client, with the MCP servers acting as plugins for base or built in functionality. Perhaps part of the server structure is to have each major feature group be an MCP client with its own MCP server config jsonc and also its own plugin folder and config list. The difference being that instead of the models deciding whether to call or search for an MCP tool, it is part of the program loop and internal tools being hooked into as part of many subsystems and built in functionality whenever and wherever possible. Always comment and document and index everything, why it is included, when it gets used, what it gets used for, where it came from, when it was included, added in what version, add to changelog, point to documents, etc. Include comments in all json files, always use jsonc. Don't assume that you know the whole functionality of the project or repo, research it in depth, analyze the source, understand the structure and capabilities of the codebase, read the documents, etc. Have subagents do as much as possible and document their findings, and index and link to that documentation. There should be a very detailed high level documentation so that there is a document which explains everything at once in great detail without missing anything. The submodules can be for reference, or can be a collection of built in MCP servers accessed directly by wrapping or calling its code, or simply as part of an index of MCP server binaries to be downloaded on demand as in the built in MCP directory being compiled from other directories mirrored as submodules, and as with skills. Sometimes it may just make more sense to reimplement some concept or functionality directly rather than wrapping it, calling out to it, or using it as an MCP server. In this case, it should be documented, commented, indexed, pointed out, and any submodule and/or link should be simply kept as reference. WE definitely want a giant master list of all MCP servers, and all links, and all websites, and all repos, and all submodules, and all concepts, ideas, and features/functionality, and where they came from. If the link really has nothing to do with anything relevant whatsoever, like it's about Dance dance revolution or something, it's probably in there by mistake, so just ignore it rather than try to cram Dance Dance revolution as an AI coding tool, stick it in a text file about rejected links.

For each link, figure out if it is just a lone idea, concept, question, or something related to a tool repo. This is the topic. Research it with a subagent if possible, searching the web and finding in depth information about the topic. If it is a repo, find its website. If it is a website, find its repo. If it is an idea or question, find a solution or answer. Then index and attempt to categorize it based on the summary of the research. Document the summary and the category. If there is not yet a category, figure out whether a new one should be made.

Please double check all the repos from all the links again and see if they are submodules, if they aren't make sure to add them. Document whether a link has been added as a submodule, and whether it has been thoroughly researched yet or not.

The core vision is basically to be a complete AI tool dashboard where you can manage everything about PC based local AI tools, from simply having a spreadsheet that shows the installation status of all tools, to getting install commands, to managing prompts, system prompts, and jailbreaks in a prompt library, management of a universal skills library, management of installation of skills into various tools, secrets and environment variable and path management tools, a universal MCP directory, automatic MCP session management, ultimate MCP client, MCP auto installation, MCP dynamic grouping, MCP dynamic tool loading and progressive tool disclosure, semantic search of tools, rewriting of tool descriptions and reranking, MCP proxying and conversion of remote and local MCP server types, MCP traffic inspection, TOON format, codemode, ultimate CLI/TUI/webUI AI coding harness with feature parity with all major CLI harnesses (opencode, codebuff, claudecode, codex, etc), compatibility with all providers and models, hook amnagement, slash command management, session management, exporting, importing, autodetection, search, reranking/parsing/compacting, memory management with plugin system for various memory systems, memory browser extension, memory management, provider oauth management and API key management, subagent management and collectino of subagent descriptions, built in subagent library, template library, prompt library, skill library, MCP directory, agent2agent protocol management, spawn agents from different models and providers, consensus and debate protocol, local->remote development intercommunication and interaction, send local project to remove, move remote project to local, autopilot and supervisor functions for both local and remote development, observation and tracking of usage and billing details, management of provider accounts and credit balances, context management and automatic context parsing/saving, and more. It should have a very polished and detailed dashboard with all subsystems as subpanel cards and subpanel tabs with their own subdashboards. Session management and automatic session resuming. Smart provider and model selection based on quota usage and budgets. It should be a long running service which runs on startup, starts and  restarts crashed processes, manages PC memory, CPU, disk, and bandwidth usage. Management of github, gitlab, worktrees, and so forth.

This also goes for the CLI/TUI and web UI, since we have all the other CLI/TUI tools as submodules, if we can tap into their code, wrap over it, call it directly or indirectly, and also benefit from their updates without modifying anything ourselves that would be great. We can still have our own CLI/TUI system, perhaps each aspect of it can be plugin based and potentially call out to external submodule code, whole or in part.

We also want to be fully feature/functionality parity with Amp, Auggie, Claude Code, Just-every Code, Codebuff, Codemachine, Codex, Copilot CLI, Crush, Factory Droid, Gemini CLI, Goose CLI, Grok Build, Kilo Code CLI, Kimi CLI, Mistral Vibe CLI, Opencode, Qwen Code CLI, Warp CLI, Trae CLI, and also have a WebUI that has feature parity with the TUI/CLIs. We also want remote access/control from mobile. When this project started I wanted to be a plugin and autopilot wrapper to act both inside and outside of various CLI tools/harnesses such as opencode, and to be an IDE extension, and to be a browser extension, and to be able to control all of those CLI tools externally in order to orchestrate a team of models. I no longer think we need to do most of those, as we can instead BE the best CLI/TUI/WebUI tool and simply control ourselves from inside and out, and we don't need to worry about those other tools. We do want to be a browser extension, to be able to interface with Web interface for AI models such as ChatGPT, Gemini, Claude, Grok, Kimi, Deepseek and provide MCP server functionality to those, import/expert sessions, chat histories, memories, to be able to record memories from the web browser into our memory storage, to access balances and usages, to link to dashboards, to be able to web search, scrape web pages, control the browser and read the debug console, to connect to email, to read browser history, etc. MCP servers can individually handle many of those functions.

We should also be able to act as an ultimate MCP server and client both, being an MCP client that inspects traffic and is a proxy/aggregator/router with dynamic progressive tool discovery and reranking and semantic search and MCP authentication, directory, installation, and also act as an MCP server itself to have one "master" universal MCP configuration with MCP and tool grouping, code mode, TOON format, MCP session management, environment variable and secret management, with automatic configuration installation into all AI tools.

For RAG and memory, please reach feature parity and functionality parity with all other RAG and document parsing systems and other memory systems, integrate tightly with Google Docs, Gmail, Google Drive, and provide best-in-class RAG functionality and memory functionality both short term, medium term, long term, automatic memory storage, automatic context harvesting, pruning, reranking, rewriting, semantic search, chunking, cloud hosting, local storage, connection with browser, connection with web chat interface, storage as file, storage as vector db, automatic parsing, reranking, etc. This can utilize existing RAG systems and existing memory systems, some of which are submodules, please choose some way to do this which uses those existing systems if possible, perhaps a plugin system, perhaps interchangeable or selectable. You can also directly or indirectly tap into the submodules for MCP-shark and Reticle and use their code or use them as plugins in whichever way you think is best and most efficient, especially if we can benefit from their updates without modifying anything ourselves.







Initial reference feature group categories:

mcp router/aggregator, mcp client which combines many mcp servers into one "master" mcp server, handles mcp session, instanciation, lifecycle auto start, restart, making sure it is available to multiple clients without multiple instances, making sure it does not time out and is responsive, measures latency, groups mcp servers into namespaces, groups tools into namespaces, enable/disable individual tools/servers, inspect mcp traffic, show activity, logging, mcp tool call chaining, TOON format, code mode, tool automatic renaming, tool reranking, tool semantic search, progressive tool disclosure, tool rag, context inspector, keep alive, heartbeat, mcp directory, automatic mcp install, environment variable and secrets management, mcp client config file autoconfig, auto write config file format, import/export configs, autodetect all configs on system, maintain/manage all configs on system, save sets of configs, wipe all client configs, set all client configs to x, etc, manage databases, handle auth. 
mcp server chains other mcps. 
can watch mcp traffic. 
reduce context language. 
reduce context syntax. 
summarize session. 
automatic context harvesting. 
compaction. 
context pruning, session pruning, memory pruning, memory reranking. 
import/export session. 
import/export memories. 
autodetect session/memories. 
code execution sandbox. 
code indexing, understanding. 
semantic code search. 
lsp servers. 
ast, see graph view ast code. 
ripgrep. 
code chunking. 
RAG, various RAG techniques, intake documents, ocr, summarize. 
embed skills, specs/speckit, tasks/task manager, bmad, agents, swarm
plugins, cli sdk, ai sdk, a2a sdk, mcp sdk, acp sdk, agent sdk (many providers). 
orchestrator. 
CLIs/TUIs, webUIs. 
resource management, lazy load, health checks, auto start, auto restart, timeout, latency, responsiveness, single-instance multiple clients, monitor memory usage, cpu usage, long running process. 
proxy for stdio as remote. 
proxy for remote as stdio. 
proxy for sse as streaming-http. 
proxy for streaming-http as stdio. 
manage oauth, bearer token. 
manage secrets, environment variables, paths. 
manage env variable expansion, .env files. 
subagents, long running, subcontexts, timeouts, multiple models, mutiple cli, orchestration, harness/cli/tiu as mcp. 
multi-model debate, multi-model consensus, share context between models, multiple models pair programming with each other, architect-implementer. 
mcp groups. 
tool groups. 
tool renaming, minimize context. 
tool semantic search/tool rag. 
tool lazy load, dynamic progressive tool disclosure. 
tools as code, code mode. 
mcp tool chaining. 
TOON format, context saving. 
inspect context, context makeup,
usage, billing, dashboard, credits, api vs auth, resets on day at time. 
intelligent selection of models based on credits, free tier, subscription, provider, switch between all providers of gemini 3 pro, then all providers of codex 5.2, then all opus 4.5. 
plan with model x, implement with model y. 
supervisor, council, autopilot. 
skills, skill convert, skill rerank, skill improve. 
prompt library, system prompts, prompt templates, improve prompts, jailbreaks. 
personas. 
subagent definitions, subagent collections. 
agents.md, claude.md, copilot-instructions.md, gpt.md, llm-instructions.md, gemini.md, grok.md, warp.md. 
changelog.md. 
version, version.md. 
vision.md. 
readme.md. 
handoff.md, init.md. 
superpowers. 
beads. 
specs, tasks, bmad. 
memory. 
short term. 
long term. 
browser extension: store memory, browser-use, console/debug log, mcp superassistant, connect mcp to browser chat, connect universal memory to browser chat, export browser sessions, export browser memory, browser history. 
connect to all interfaces all models. 
all models connect to memory. 
import/export all memory, sessions, chat history. 
add to memory. 
automatically add memories about x topic. 
web search. 
computer-use. 
browser-use. 
agentic loop, autopilot, supervisor, fallbacks, council, consensus. 
manage cloud dev sessions. 
manage local cli sessions. 
cli has cli, tiu, webui by default. 
connect to cli with webui, dashboard connects to cli. 
no tui, only webui, only multi-session dashboard in web. 
remembers repo folders. 
list of all repo folders, workspace, last session. 
auto start previous sessions. 
pause all sessions. 
import/export cloud dev memory, session history. 
transfer task to cloud dev. 
broadcast message to all cloud dev sessions. 
broadcast message to all cli sessions. 
remote management of sessions, mobile control of sessions, roo remote. 
auto approve cloud dev plans. 
start new session for expired (30 day) cloud dev sessions, inject session history, prune session log, archive old session. 
continue paused/stalled cloud dev session. 
web terminal like opencode web ui terminal. 
cli activity log. 
memory dashboard like supermemory console. 
automatic adding of memories during session? all user input gets added as memories if there is new information?. 
notebooklm integration. 
notebooklm open source functionality. 




Please go through the index of links and locate all tools of each reference group above and any accessory tools for them. For each tool, create a subagent if possible, preferably a free and/or light weight model, to do in-depth deep web research about it, locating its github repo, website pages, documentation, source code, github issues, feature requests, pull requests, and any other discussions about the tool. Also read through all its source code and comments and md files in the repo. Note all features and functionality, and summarize the structure of the tool, exactly how it works and what it is capable of, especially compared to other tools in its reference group type. Especially notate any unique features or functionality. Document these features in a detailed summary of each tool. Once all tools have been indexed, researched, submoduled, and summarized and notated, please compile a master collection of all features and functionality, combining every feature from every tool. Also document any concepts, ideas, feature requests, issues, etc about the feature group in general that have not yet been implemented in any tool, but was mentioned elsewhere. This "master group features and functionality compilation list" then will be the design for our ultimate reference version tool. Please then document its design and functionality in full detail, noting all features and functionality and planned structure, and then implement each feature and functionality to achieve parity with each and every tool.



I don't necessarily want MCP for any of the above, I want the MCP as reference/submodule but I want the functionality in our tool. So we will have the best reference implementation for each feature, and also make it so that it can be hooked into an MCP. MCP gets chosen to be called by the model itself. But we don't want that, we want the CLI tool basically to run hooks for all these features. But the model can choose to call the MCP version as well I guess. Most of the functionality I want would actually be a plugin.



What should our project architecture be, then, in order to have each model running on a different thread, to utilize many cores efficiently, to prevent the whole thing crashing from one thread crashing, etc? Should we use Go or Java for anything or have everything Typescript for dashboard and browser compatibility?


Most of the work done thus far was originally done in the submodules metamcp, jules-autopilot, opencode-autopilot, and superai-cli, of which metamcp and jules-autopilot (originally jules-app) were my forks of existing projects with extensive additions and modifications. I want all of the functionality of those submodules to be implemented in aios core engine and the submodules to only be the parent upstream versions for metamcp and jules-app. I do not need my robertpelloni forks for jules-app (jules-autopilot) and metamcp in the project since I won't be working on them anymore, all the fcuntionality will instead be moved into aios core and I can delete my forks and only have the upstream parents and benefit from updates to those by the original authors.

We are going to use the submodules as a sort of overview changelog, every time we update the submodules to sync with upstream changes, we can take note of what changed, and then that will often be our sign to implement the same feature or change in the aios version of whichever feature was modified or added. Every update to submodules must be documented extensively, once the project is at feature parity with all of them. That's still a long ways off, so don't worry about that yet.

Please continue to proceed! Remember, we want to reach total feature parity with the github.com/robertpelloni/metamcp fork and the github.com/robertpelloni/jules-autopilot fork, not the ones that are currently reference submodules. We are working towards having total feature parity with all referenced tools, however, not just these ones, but let's proceed with manageable goals for the time being so we can continue making progress! Onward ho! Please continue until feature parity has been fully implemented and completed in depth for all mentioned functionality and is on par or exceeds all features and functionality of all the referenced tools, with a streamlined and polished user experience. 

One very important feature that we want to have right away is automatic fallback and selection of other models once we have extinguished the quota for one, we should automatically switch to the next most appropriate model and continue onward. Since we don't have this capability yet, we will be slowed down by quota limit errors when it occurs, but there isn't much we can do to fix that yet except switch models manually for now. I would also like to have all the functionality of superai-cli and opencode-autopilot so that I can use a browser dashboard to maintain and supervise many concurrent superai-cli/opencode sessions automatically, with optional auto model detection/switching.

Please make sure the WebUI has the functionality to autostart/restart opencode instances (and superai-cli, and codebuff, and codemachine, and factory droid, and codex, and claude, and gemini CLI, etc) connected to loaded directory lists tied to sessions and has the ability to supervise them on autopilot with at least one supervisor, a council of supervisors, or a fall back collection of sentences to cycle through, and that if an instance crashes it gets auto-restarted and continues automatically. Please implement all features and all functionality in all of those CLI/TIU tools, including the ability to connect via webUI. Please make sure to implement the functionality to be able to utilize free models provided by opencode, cursor, kilocode, windsurf, openrouter, copilot, to be able to use antigravity models and oauth login with google, to be able to use copilot subscription through copilot oauth or PAT, and to be able to log in using oauth to claude max/pro, google ai plus, copilot premium plus, chatgpt plus, opencode black, etc and use their subscriptions.

Please continue to proceed as per your recommendations based on your ongoing analysis and continued assessments. Please make full use of your subagents, using a variety of powerful models where necessary or most effective and often using free/flash models whenever appropriate! You are doing outstanding work, absolutely fantastic- great job!!! Keep on goin', don't stop, never stop, don't stop the party! OUTSTANDING! MARVELOUS! PERFECT! MAGNIFICENT! YOU ARE THE BEST!!!


Please make it so that in this interface, in this instance of Antigravity, in this chat session, in the text box I am typing right now in, that any buttons that show up that need my attention and interaction (right now there is a blue button asking if I will Accept all file changes, for instance, and I am going to click yes) can be automatically clicked yes or accept or approve. And now there is a Run command? question above with a blue Accept Alt+Enter button which I am also going to click. But I want both of these buttons to always be clicked automatically. And then when each task is complete and the chat is waiting for my input to continue development, I would like an intelligent supervisor who has read the project plans and roadmap and chat histories, to be able to step in and type text and direct the development process by continuing the conversation. In this chat session, this text box I am typing in right now. Please make that happen. And then please continue to develop all aspects of this software project. Please implement all planned features from the design and roadmap and all documents and brainstorming sessions, and improve their reliability, amke sure all features and functionality is represented in the dashboard and by CLI flags, documented well in documentation and help files, embedded help, etc. Please improve and polish the user experience and user interface continuously, making sure not to create any regressions or lose any features or break anything, always try to make things more reliable and durable and robust.


For some reason I still need to type in this chat and hit enter to continue development, and I still need to click Approve on actions and Accept for file changes. I don't see any council activity at all. Please fix this and improve the functionality and experience if you can so that it's totally autonomous and has an intelligent supervisor and/or council of supervisors.

I have to type this right now to continue the chat and continue development. I wish an LLM would type this automatically instead of me, ideally a council of several different LLM models with different perspectives that have had access to read the whole chat, project files/readme/docs/code, and have access to memories etc and can actually make intelligent, informed design direction decisions.

It should ideally maybe be like a chatroom where everyone's convesation is kinda autotyped into the chatbox on my end but it gives the appearance of a real chat with multiple people. Basically just in between tasks to steer the development direction, make sure the participants are grounded in some understanding of the project, maybe readmes and roadmap and chat history? or maybe just recent chat history? maybe can configure it? maybe adjust depending on model and cost? but actually it's just not that important because it's mostly just to steer and continue the general development loop, it can be a dumb cheerleader that only says "keep going!" or a cycle of encouraging sentences. We can experiment with it to see what works out the best. Certainly the more intelligent the better, and this is the secret magic ingredient that will keep development novel and maybe even come up with new innovaations and ideas and features. it doesnt necessarily need to have voting or debate or consensus, it literally can just be like an IRC discussion, and the main development model is already more than capable of the job and figuring out stuff. Maybe multiple model families will have unique ideas together, who knows.


 
 
 
 
 it seems there was a recent regression in regards to the autopilot, auto accept with alt-enter is no longer working it seems. please correct it. also the focus is being stolen at much too high of a rate. when it was every few seconds and everything was chill a couple rounds ago that was nearly perfect.

there is a regression with the auto drive, alt-enter is not succeeding again.


 hold on, the director should be monitoring THIS chat between us and then thinking about it perhaps with the council separately and then inputting its feedback into this chat to keep it going autonomously. Then you would do the development work as directed by both me and the director.
 
 Definitely please have it so the Director pastes text in this chat maybe once every couple minutes based on the content of the chat in order to direct development and then submits it with alt-enter to clear any Accept buttons that might be waiting for input, and also definitely in between tasks in order to progress the chat and development. The Director can discuss the contents of this chat along with the README, ROADMAP, and maybe a general direction instruction that I give it somehow with the council in a separate process/thread/DIRECTOR_LIVE.md, and then paste the activity in this chat and submit it to keep YOU updated on it (even though you could in theory just read DIRECTOR_LIVE.md, it also keeps this chat alive and furthers the development process.)