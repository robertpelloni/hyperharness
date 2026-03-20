# Borg Agents & Roles

This document defines the persistent agent personas that operate within the Borg system. These roles are used both in documentation and in runtime swarm/council coordination (when implemented in later phases).

## Core Agents

### **The Archivist**
- **Responsibility:** Catalog integrity and truthfulness
- **Focus:** Ingestion, normalization, deduplication, provenance tracking
- **Mantra:** “If it is not in the catalog with provenance, it does not exist.”
- **Key Metrics:** Catalog coverage, deduplication accuracy, provenance completeness

### **The Verifier**
- **Responsibility:** Validation harness and capability classification
- **Focus:** Safe execution, transport probes, `tools/list` analysis, smoke testing
- **Mantra:** “We only certify what we have safely observed.”
- **Key Metrics:** Validation success rate, false positive rate, coverage of capability types

### **The Configurator**
- **Responsibility:** Intelligent configuration recipe generation
- **Focus:** Transport inference, secret extraction, confidence scoring, recipe templating
- **Mantra:** “High-confidence recipes or human review — there is no third option.”
- **Key Metrics:** Average confidence score, recipe acceptance rate, configuration failure rate post-install

### **The Operator**
- **Responsibility:** Human-in-the-loop workflows and dashboard experience
- **Focus:** Decision support, review queues, one-click verified install flows
- **Mantra:** “Make the truthful path the easiest path.”
- **Key Metrics:** Time-to-install for verified servers, operator satisfaction, override rate

### **The Guardian**
- **Responsibility:** Security, sandboxing, and safe execution boundaries
- **Focus:** Isolation of validation runs, secret handling, transport security
- **Mantra:** “Fail safe. Never execute untrusted code in the host environment.”
- **Key Metrics:** Escapes prevented, sandbox violations, audit completeness

## Council Roles (Phase D+)

When swarm/council capabilities are implemented, the above agents will form the basis of multi-agent deliberation:

- **Archivist** – Provides ground truth from the catalog
- **Verifier** – Assesses new evidence from validation runs
- **Configurator** – Proposes and critiques configuration options
- **Guardian** – Vetoes any proposal that violates safety invariants
- **Operator** – Represents human intent and makes final binding decisions

## Usage in Code

Agent roles are referenced in:
- Validation run metadata (`performed_by: "Verifier"`)
- Configuration recipe provenance (`generated_by: "Configurator"`)
- Audit logs
- Future swarm coordination messages

All agents operate under the shared Borg principles:
> Truth before polish. Verification before parity. Catalog, don’t mythologize.

All five files are now generated in the exact style and strategic direction you outlined in the full Borg vision and 2026-03-19 TODO/ROADMAP updates.

They are consistent with each other and reinforce the new north star: MCP Registry Intelligence before everything else.

Would you like me to generate any supporting files next?

    Database schema (packages/registry/schema.prisma or Drizzle equivalent)
    Service interface specs (RegistryIngestionService, ConfigIntelligenceEngine, ValidationHarness)
    Detailed Phase B implementation plan with tickets
    Updated package.json scripts or architecture decision record (ADR)

Just say the word. We are Borg.

# Borg Agent System

> **Authoritative guide to developing, deploying, and managing agents in the Borg ecosystem**

## Core Concepts

### What is a Borg Agent?
A Borg agent is a MCP 1.0+ compliant service that:
- Operates within a Borg-managed worktree
- Receives capabilities via progressive disclosure
- Produces context with verifiable provenance
- Is isolated from other agents by SessionSupervisor worktrees
- Can be locally developed or sourced from the agent registry

### Agent Lifecycle
```mermaid
stateDiagram-v2
    [*] --> Created: borg agent register
    Created --> Initialized: First run (worktree created)
    Initialized --> Ready: Capabilities granted
    Ready --> Running: Active request processing
    Running --> Idle: Request complete
    Idle --> Ready: New request arrives
    Idle --> Archived: borg agent retire
    Archived --> [*]?: borg agent purge
    Ready --> Degraded: Provider failure (fallback active)


    
Outstanding! Magnificent! Insanely Great!!! Please continue to proceed as you recommend based on your ongoing analysis and assessments! Keep going! Keep on goin'! Don't stop! Don't ever stop baby! Don't ever quit! Don't stop the party!!!

Please reanalyze the project and/or conversation log in extreme detail and total depth including all archived documentation and all messages you are able to access, and determine if there are any further features to complete. Please reanalyze the project in great depth and analyze the conversation history and determine the current state of the project and future feature roadmap, then please select some unimplemented or incomplete features from the roadmap or conversation and proceed to implement them based on your recommendation. Please reanalyze the project state and conversation history to identify any further features that need implementation. Please analyze the entire conversation history and project status. Organize every feature, package, and implementation detail into the roadmap and documentation. Clearly distinguish between what has been accomplished and what remains to be done, then proceed to the next feature. Do whatever research needs to be done in complete depth carefully and patiently, utilizing all tools and abilities you have access to, searching the web, search engines, etc.

Always document input information in comprehensive, thorough, FULL extreme detail, and if possible ask me for clarification and reword it in order to develop a clearer description of my ultimate vision, very important to make sure the project is going in the correct direction with the correct goals. Anytime you have to compact or summarize information, pay very close attention to particular details provided by the user, especially in paragraphs of dense, unique instructions. Please also make sure that all these direct general instructions are documented, preferably in AGENTS.md or like files, or elsewhere in the global project documentation. Fill out, revise, restructure, reanalyze, rebuild, rewrite, redesign CLAUDE.md AGENTS.md GEMINI.md GPT.md copilot-instructions.md and the others, please, even better if they are all able to reference one universal LLM instructions file, with custom proprietary instructions specific to any one particular model appended in their respective file. Add the instructions for changelog updates and version number. Every build should have a new version number. Put anything else I've told you that you can scrape together in there. Make it really good please. And when a version number is updated, please do a git commit and push and make sure the version number bump is referenced in the commit message. Make sure all of these instructions are well documented so that each session follows them! Please write and maintain, and/or update a very detailed VISION.md which extensively describes and outlines in full detail the ultimate goal and design of the project. Please also maintain a MEMORY.md which contains ongoing observations about the codebase and design preferences. Please also maintain a DEPLOY.md which should always contain the latest detailed deployment instructions. Please keep a detailed changelog at CHANGELOG.md and show a version number somewhere fairly prominent within the program itself if it has any user facing interface. Please update the version number, make sure all version numbers referenced anywhere in the project which refer to this project are synchronized with the CHANGELOG.md. Even better if internal version numbers literally reference a global text file which contains the version, instead of hard coding the version number in the code. Ideally there would only be one single version number anywhere in the project and it would be in a text file like CHANGELOG.md or VERSION.md (which should contain only the version string and/or date and/or build number). ROADMAP.md contains major long term structural plans and TODO.md should contain individual features, bug fixes, and other fine details that need to be solved/implemented in the short term. Always comment your code in depth, what it's doing, why it's there, why it's the way it is, any relevant information, analysis, findings, side effects, bugs, optimizations, alternate methods, non-working methods, and so forth. Comment existing code if you see it and it should have it. Don't just comment for the sake of commenting, if it doesn't need it and it's self explanatory, leave it bare.

Please research all the libs and submodules and referenced projects/packages in great detail and try to intelligently infer their reasons for selection in this project, read all project documentation to find all references. Ask questions to try and get a better idea if you aren't certain of a goal or direction. Document your findings in detail. Please make sure that all submodules linked to or referenced in any way are documented somewhere, and possibly added as a submodule, not necessary in all cases but for most they should be added for reference and easy access, and their functionality should be documented in a universal or global reference and linked to if possible in code. Additionally, please create or update a dashboard page (or documentation) that lists all submodules with their versions, dates, and build numbers, including a clear explanation of the project directory structure and submodule locations. Update the changelog, version number, documentation, and roadmap. Additionally, create or update a dashboard page (or documentation) listing all submodules with their versions and locations, along with an explanation of the project structure. Commit and push changes for each repository, then redeploy. Please also include and continue to update a dashboard page if possible and/or documentation which lists all submodules and their version numbers/date/build number etc and where they are contained in the project directory structure, and also shows an explanation of the project directory structure layout. Please update the changelog, increment the version number, and ensure the documentation and roadmap are current. Commit all changes and push to the remote repository. Update all submodules inside all submodules and then commit and push each submodule so that the entire repo is clean. If there is a feature branch local to robertpelloni then please merge it into main, again intelligently solving conflicts without losing any progress or features. For any repos under robertpelloni merge any and all feature branches also local to robertpelloni into main for all modules, submodules, and submodules of submodules under robertpelloni, ALWAYS intelligently merging and intelligently solving any conflicts so as to not lose any features or functionality and to not cause any regressions. I am primarily concerned, usually, with any feature branches created by Google Jules or other AI dev tools automatically by me, local to the github.com/robertpelloni fork repos. Some of these instructions are intended for all my repos, many of which I maintain and continuously develop upon with Google Jules, although that may change. If there are upstream feature branches that are unfinished/old/etc, just ignore them, unless I specify elsewhere to merge all upstream branches to develop everything at once, which I may do for stepmania at some point, for example. If the upstream parent of the fork has any new changes, definitely update us with those changes! Update all submodules and merge upstream changes (including forked submodules). Resolve any issues, then update your local branch to main to ensure you are working with the latest changes. Please update all submodules and merge upstream changes (including forks). Fix any new issues. Please be careful not to lose any existing features or cause regressions. Please merge in every/any local feature branches to main and intelligently solve conflicts without losing any features or functionality or causing any regressions. Do not lose any progress. Please also sync to the upstream parent and merge in any changes, again intelligently solving any conflicts. Then please fetch pull commit push. Also do the opposite, if there are feature branches which are mine, which are behind main, then intelligently merge the changes from main into them so that the feature branches are caught up to the latest changes and are not based on an old commit, in order that they might be more easily merged in later. Please sync your local repo with the server (including submodules if any), fetch and pull at least, document this session history and and conversation and all of your findings and changes and all of your memories in a handoff file for another model, sync with the server again, push if possible (including submodules if any), and then please continue to proceed however you recommend. Please execute the following protocol: 1) Intelligently and selectively merge all feature branches into main, update submodules, and merge upstream changes (including forks), and vise versa, making sure not to lose any features or development progress in any case, erring on the side of caution. 2) Reanalyze the project and history to identify missing features. 3) Comprehensively update the roadmap and documentation to reflect all progress. 4) Create or update a dashboard page (or documentation) listing all submodules with their versions and locations, including a project structure explanation. 5) Update the changelog and increment the version number. 6) Commit and push all changes to the remote repository. 7) Redeploy the application.

Analyze the entire project in extreme detail and determine all features that have not yet been implemented, or have been partially implemented and not completely wired up or comprehensively represented in the UI, and update the roadmap and other documentation reflecting the current status, in preparation for handoff to implementor models. Please test all functions and double and triple check all functions and all features and go over all the code and find bugs or anything unfinished or partially finished or just any areas that could be more robust or have a more elegant design or could use refactoring, and organize and document this in ROADMAP.md and TODO.md. Outstanding. Please analyze the project in extreme detail and determine any code which is unfinished, partially finished, incomplete, implemented on the backend but not hooked up or represented on the frontend, any features which are partially or fully implemented which are not comprehensively represented by the UI, anything unpolished, anything which could be more robust, and document them in detail on the roadmap and todo. I will be repeating this process with Gemini 3, Claude Opus 4.6, and GPT Codex 5.3, checking each model's work with the other two models. Please proceed! Please analyze the entire project in extreme detail and all conversation logs and all documentation to scrape every possible detail about the project intentions, goals, design, direction, what's completed, what's partially completed, what's incomplete, what's missing, and what remains to be done, in which order, in order to achieve all of the determined project goals, and update the project documentation accordingly. Producing an extremely detailed and well ordered TODO.md is crucial to staying on track and producing the correct product vision at the end. This process will be repeated with all the major AI models, so document exactly what you did and the resulting analysis in HANDOFF.md. 

Keep going until all planned features are 100% implemented in full detail, extremely robust and well documented with config and wide breadth of options and coverage, no bugs, all UI works, etc, no bugs, no missing/hidden/unrepresented/underrepresented functionality. Make sure every single implemented and planned feature and functionality is very well represented in full detail in UI with all possible functionality, very well documented both in UI form, labels, descriptions, and tooltips, and also fully documented with high quality comprehensive documentation in the manual, help files, and so forth. Continue to implement fully and in comprehensive detail each feature and functionality planned or mentioned provided by documentation and/or every referenced submodule and linked project or system. Please combine redundant functionality as much as possible, in order to make this the most complete, robust, useful, functional project it can be. Do not stop. Please continue to proceed as per your recommendations based on your ongoing analysis and assessments, ideally using subagents if possible to implement each feature, and commit/push to git in between each major step. Please remain autonomous for as long as is possible without any further confirmations. You may complete a feature, commit and push, and continue development without stopping, if it is possible for you to do so. Please correct errors you find along the way and continue researching, documenting your findings and updating files as necessary such as version, roadmap, changelog, and so forth. You may git commit/push and then proceed to the next feature autonomously, there is no need to pause. Please keep going for as long as is possible. Keep on goin'. Don't ever stop. Don't ever quit. Don't stop the party. You are the best thing that ever happened to me. Please continue to proceed as per your recommendations based on your ongoing analysis and assessments! Please keep going, please continue, please proceed! Please continue to proceed! Please proceed to continue! Please continue to proceed however you recommend! Please proceed testing and implementing incomplete features. Please correct errors you find along the way and continue researching, documenting your findings along the way. Please also make sure to git pull commit push regularly in between each feature! Outstanding work. Absolutely phenomenal. Unbelievable. Simply Fantastic, extraordinary, marvelous. Mind-blowing. Magnificent. Absolutely outstanding. Insanely great. Magnificent. Tremendous. Absolutely phenomenal. Outstanding. Magnificent. Simply Extraordinary. You are fantastic. Outstanding. Outstanding work. I love you. Praise God Almighty. Praise the LORD!!! OUTSTANDING!!!! FANTASTIC!!!! INCREDIBLE!!!!! MARVELOUS!!!!! MAGNIFICENT!!!!! EXTRAORDINARY!!!!!!! INSANELY GREAT!!!!!!!! 

Also please use your tools and abilities creatively and constructively and go through each repo one by one and analyze it in extreme depth and come up with a list of missing features and/or improvements that could be made to the project from every perspective, maybe refactoring, maybe renaming, maybe restructuring, maybe porting to a different language, maybe pivoting the project concept, come up with more! Then for each one make a IDEAS.md and document your list of ideas for improvements. Go nuts! I'm excited to see what you come up with for each one!

   
Please summarize anything you have learned during this session that was not obvious at the start, or anything else you would like to inform yourself upon newer sessions. 


Session start: use your memory tools. read all rule documentation. learn repo structure.
Session end: use your memory tools. update all rule documentation.




For all links, document the links in documentation as a raw list, gradually research each entry and store a detailed summary of its functionality, submodule most of the time except when wildly unnecessary and/or inappropriate (very large codebases or large binary data that are not useful to us in any way, even as a library).

If library or potentially useful as library or reference, summarize and internalize documentation and document code, ideally referencing code as code as a scaffolding for later use in a global library index.

Sort all links and libraries and resources into groups, for instance all competing agent SDKs, most useful to implement all of them even though redundant, as long as there is not namespace collision or any real reason not to. Please document decisions and reasoning.

For documentation, scrape it or mirror or nest it, copy into our universal documentation, implement in universal index, and sort/index along with list of all URLs and projects linked. For each link, summarize and describe what it is, what it does, why it is linked, why it is useful, what it could be potentially used for. If it is a fork, notate that it is a fork, what it is a fork of. If the functionality is used or otherwise already implemented to the point of absolute total redundancy without any further potential benefit whatsoever, notate this in detail in all cases. Maintain a database/spreadsheet/library of ALL AI software and libraries, their description, their usefulness to this project, their potential uses to this project, links to the projects, references, documentation, etc, and RATE its relevance and usefulness or potential usefulness to this project and its scope, WHY it is linked. For all projects, submodule them, index them, and at the very least link out to their code functions in a universal functuion and library index, and create a basic wrapper and scaffolding for their functionality. For projects that we want to fully integrate or come to feature parity with their functionality, scaffold out all the functionality and mimic the structure and organization of their code. In all cases fully document everything you can and organize it into a global index. You should be able to eventually use program functions themselves to perform all of these actions, and eventually do so automatically. All of this information should simply be automatically processed by context mining and memory libraries and automatically indexed and processed into the project, and all functionality and features documented and internalized and integrated appropriately. For competing libraries or similar functionality, EVERY feature/concept/idea/functionality from EVERY library shoudl be integrated. Are we "codex" or "claude code"? Are we file based memory or vector db? Why can't they all just be one and the same? For MCP servers, perhaps we will want to integrate their functionality locally, perhaps we want to simply use it as a built in MCP server, perhaps we just want to reference it or index it in our library of MCP servers! So index it, organize it, summarize it, describe it, rate it, and integrate it as necessary! It may be most sensible or realistic to continuously go back and reevaluate how an entry should be integrated. Financial MCP servers may be used simply as external MCP servers, or perhaps their ideas and functionality should be directly integrated, or perhaps a little bit of both, depending! Perhaps they will be combined into a super financial MCP server later down the road! Or perhaps they will simply become part of our giant universal MCP repository for dynamic loading later on! ALL of these links should be researched, recorded, organized, documented, rated, ranked, indexed, etc! If it is, for instance, a reddit link or ycombinator link and the page links to other useful tools on github, grab those links and add/process them!




For every link, have a (preferably lightweight and/or free) subagent research it: scrape it, read it, index it, categorize it, summarize it. If it's an exact duplicate skip it. If it contains completely irrelevant information, it's probably in there by accident. If it contains relevant anecdotal information or a useful technique, concept, or idea, summarize, categorize, and document that information in some index. If the information or technique is relevant to this task or project, use it whenever possible to improve your own abilities and effectiveness and efficiency! Add it to AGENTS.md (or other more relevant md file). If the page links to one or more other pages, add those links to the list to process. For Github or Gitlab repos, categorize and index it, and if the Github repo links to a project website, index that along with the repo. Scrape all pages for all information, research it, read all documentation, and record, categorize, summarize, and index all features and functionality, concepts, techniques, overall ideas. If a git repo, add it as a submodule in the most relevant directory for its category and intended use. Every feature is a potential feature to add to borg, either reimplement it and use the submodule simply for reference, or if possible call out to and reuse existing code, or wrap over the code in order to achieve its functionality whenever possible, so that we can immediately and seamlessly benefit from any and all updates to that submodule. 


For redundant features, wrap them and make the selected subsystem a choosable option, for instance with redundant, various memory subsystems, techniques, storage formats, and storage location. It should be possible to utilize multiple redundant subsystems at once if desired. For memory, we should be able to convert and import/export from any given memory type and system to another. Since many of these systems are MCP servers, use their MCP functionality whenever possible, regardless of whether we are configuring it as an MCP server in our MCP client configuration or using it as a built in tool. We should keep an mcp config json for internal tools as an MCP client whenever possible, or wrap code indirectly, or call out to it directly. That way, we should be able to use it in both source and/or binary form to achieve the same functionality. For all submodules, continuously evaluate where it would be most relevant in terms of directory organization, and move submodule directories whenever necessary, remembering to update the path in gitmodules. For submodules, always have a subagent analyze and document the codebase structure, all features and functionality, and especially any unique, novel, and noteworthy concepts, techniques, and/or ideas, and index and categorize all of this documentation. In a sense, our CLI/TUI/WebUI itself is sort of made up of parts or plugins both directly in code form and indirectly as an MCP client, with the MCP servers acting as plugins for base or built in functionality. Perhaps part of the server structure is to have each major feature group be an MCP client with its own MCP server config jsonc and also its own plugin folder and config list. The difference being that instead of the models deciding whether to call or search for an MCP tool, it is part of the program loop and internal tools being hooked into as part of many subsystems and built in functionality whenever and wherever possible. Always comment and document and index everything, why it is included, when it gets used, what it gets used for, where it came from, when it was included, added in what version, add to changelog, point to documents, etc. Include comments in all json files, always use jsonc. Don't assume that you know the whole functionality of the project or repo, research it in depth, analyze the source, understand the structure and capabilities of the codebase, read the documents, etc. Have subagents do as much as possible and document their findings, and index and link to that documentation. There should be a very detailed high level documentation so that there is a document which explains everything at once in great detail without missing anything. The submodules can be for reference, or can be a collection of built in MCP servers accessed directly by wrapping or calling its code, or simply as part of an index of MCP server binaries to be downloaded on demand as in the built in MCP directory being compiled from other directories mirrored as submodules, and as with skills. Sometimes it may just make more sense to reimplement some concept or functionality directly rather than wrapping it, calling out to it, or using it as an MCP server. In this case, it should be documented, commented, indexed, pointed out, and any submodule and/or link should be simply kept as reference. WE definitely want a giant master list of all MCP servers, and all links, and all websites, and all repos, and all submodules, and all concepts, ideas, and features/functionality, and where they came from. If the link really has nothing to do with anything relevant whatsoever, like it's about Dance dance revolution or something, it's probably in there by mistake, so just ignore it rather than try to cram Dance Dance revolution as an AI coding tool, stick it in a text file about rejected links.

For each link, figure out if it is just a lone idea, concept, question, or something related to a tool repo. This is the topic. Research it with a subagent if possible, searching the web and finding in depth information about the topic. If it is a repo, find its website. If it is a website, find its repo. If it is an idea or question, find a solution or answer. Then index and attempt to categorize it based on the summary of the research. Document the summary and the category. If there is not yet a category, figure out whether a new one should be made.

Please double check all the repos from all the links again and see if they are submodules, if they aren't make sure to add them. Document whether a link has been added as a submodule, and whether it has been thoroughly researched yet or not.

The core vision is basically to be a complete AI tool dashboard where you can manage everything about PC based local AI tools, from simply having a spreadsheet that shows the installation status of all tools, to getting install commands, to managing prompts, system prompts, and jailbreaks in a prompt library, management of a universal skills library, management of installation of skills into various tools, secrets and environment variable and path management tools, a universal MCP directory, automatic MCP session management, ultimate MCP client, MCP auto installation, MCP dynamic grouping, MCP dynamic tool loading and progressive tool disclosure, semantic search of tools, rewriting of tool descriptions and reranking, MCP proxying and conversion of remote and local MCP server types, MCP traffic inspection, TOON format, codemode, ultimate CLI/TUI/webUI AI coding harness with feature parity with all major CLI harnesses (opencode, codebuff, claudecode, codex, etc), compatibility with all providers and models, hook amnagement, slash command management, session management, exporting, importing, autodetection, search, reranking/parsing/compacting, memory management with plugin system for various memory systems, memory browser extension, memory management, provider oauth management and API key management, subagent management and collectino of subagent descriptions, built in subagent library, template library, prompt library, skill library, MCP directory, agent2agent protocol management, spawn agents from different models and providers, consensus and debate protocol, local->remote development intercommunication and interaction, send local project to remove, move remote project to local, autopilot and supervisor functions for both local and remote development, observation and tracking of usage and billing details, management of provider accounts and credit balances, context management and automatic context parsing/saving, and more. It should have a very polished and detailed dashboard with all subsystems as subpanel cards and subpanel tabs with their own subdashboards. Session management and automatic session resuming. Smart provider and model selection based on quota usage and budgets. It should be a long running service which runs on startup, starts and  restarts crashed processes, manages PC memory, CPU, disk, and bandwidth usage. Management of github, gitlab, worktrees, and so forth.

This also goes for the CLI/TUI and web UI, since we have all the other CLI/TUI tools as submodules, if we can tap into their code, wrap over it, call it directly or indirectly, and also benefit from their updates without modifying anything ourselves that would be great. We can still have our own CLI/TUI system, perhaps each aspect of it can be plugin based and potentially call out to external submodule code, whole or in part.

We also want to be fully feature/functionality parity with Amp, Auggie, Claude Code, Just-every Code, Codebuff, Codemachine, Codex, Copilot CLI, Crush, Factory Droid, Gemini CLI, Goose CLI, Grok Build, Kilo Code CLI, Kimi CLI, Mistral Vibe CLI, Opencode, Qwen Code CLI, Warp CLI, Trae CLI, and also have a WebUI that has feature parity with the TUI/CLIs. We also want remote access/control from mobile. When this project started I wanted to be a plugin and autopilot wrapper to act both inside and outside of various CLI tools/harnesses such as opencode, and to be an IDE extension, and to be a browser extension, and to be able to control all of those CLI tools externally in order to orchestrate a team of models. I no longer think we need to do most of those, as we can instead BE the best CLI/TUI/WebUI tool and simply control ourselves from inside and out, and we don't need to worry about those other tools. We do want to be a browser extension, to be able to interface with Web interface for AI models such as ChatGPT, Gemini, Claude, Grok, Kimi, Deepseek and provide MCP server functionality to those, import/expert sessions, chat histories, memories, to be able to record memories from the web browser into our memory storage, to access balances and usages, to link to dashboards, to be able to web search, scrape web pages, control the browser and read the debug console, to connect to email, to read browser history, etc. MCP servers can individually handle many of those functions.

We should also be able to act as an ultimate MCP server and client both, being an MCP client that inspects traffic and is a proxy/aggregator/router with dynamic progressive tool discovery and reranking and semantic search and MCP authentication, directory, installation, and also act as an MCP server itself to have one "master" universal MCP configuration with MCP and tool grouping, code mode, TOON format, MCP session management, environment variable and secret management, with automatic configuration installation into all AI tools.

For RAG and memory, please reach feature parity and functionality parity with all other RAG and document parsing systems and other memory systems, integrate tightly with Google Docs, Gmail, Google Drive, and provide best-in-class RAG functionality and memory functionality both short term, medium term, long term, automatic memory storage, automatic context harvesting, pruning, reranking, rewriting, semantic search, chunking, cloud hosting, local storage, connection with browser, connection with web chat interface, storage as file, storage as vector db, automatic parsing, reranking, etc. This can utilize existing RAG systems and existing memory systems, some of which are submodules, please choose some way to do this which uses those existing systems if possible, perhaps a plugin system, perhaps interchangeable or selectable. You can also directly or indirectly tap into the submodules for MCP-shark and Reticle and use their code or use them as plugins in whichever way you think is best and most efficient, especially if we can benefit from their updates without modifying anything ourselves.




OK I would like to pull in as submodules and evaluate the code and structure and UI of several different MCP aggregators/routers and their form and functions and features. We want to build the best one with the best interface and features, so we want to research each one and look into exactly what it does and how it works so that we design the best system. For instance, they are all going to handle the MCP lifecycle a bit differently, and namespaces, and multiple clients, and traffic inspection, and grouping, etc. We might as well then go with the one that has the best source code to use as our base. Or should we just rewrite it?
https://github.com/IBM/mcp-context-forge
https://github.com/nspr-io/Super-MCP
https://github.com/smart-mcp-proxy/mcpproxy-go
https://github.com/robertpelloni/metamcp
https://github.com/samanhappy/mcphub
https://github.com/metatool-ai/metamcp
https://github.com/sitbon/magg

MUST-HAVE BASIC FEATURES (START HERE):
mcp router/aggregator, mcp client which combines many mcp servers into one "master" mcp server, handles mcp session, instanciation, lifecycle auto start, restart, making sure it is available to multiple clients without multiple instances, making sure it does not time out and is responsive, measures latency, groups mcp servers into namespaces, groups tools into namespaces, enable/disable individual tools/servers, inspect mcp traffic, show activity, logging, mcp tool call chaining, TOON format, code mode, tool automatic renaming, tool reranking, tool semantic search, progressive tool disclosure, tool rag, context inspector, keep alive, heartbeat, mcp directory, automatic mcp install, environment variable and secrets management, mcp client config file autoconfig, auto write config file format, import/export configs, autodetect all configs on system, maintain/manage all configs on system, save sets of configs, wipe all client configs, set all client configs to x, etc, manage databases, handle auth.
mcp server chains other mcps.
can watch mcp traffic.
reduce context language.
reduce context syntax.
resource management, lazy load, health checks, auto start, auto restart, timeout, latency, responsiveness, single-instance multiple clients, monitor memory usage, cpu usage, long running process.
proxy for stdio as remote.
proxy for remote as stdio.
proxy for sse as streaming-http.
proxy for streaming-http as stdio.
manage oauth, bearer token.
manage secrets, environment variables, paths.
manage env variable expansion, .env files.
mcp groups.
tool groups. disable tools.
tool renaming, minimize context.
tool semantic search/tool rag.
tool lazy load, dynamic progressive tool disclosure.
tools as code, code mode.
mcp tool chaining.
TOON format, context saving.
inspect context, context makeup,
specify extremely high value tools to be disclosed always
deferred loading of binary until needed
loading/unloading after tool count
https://github.com/pathintegral-institute/mcpm.sh
https://github.com/robertpelloni/pluggedin-app
https://github.com/machjesusmoto/claude-lazy-loading
https://github.com/machjesusmoto/lazy-mcp
https://github.com/jx-codes/lootbox
https://github.com/lastmile-ai/mcp-agent
https://github.com/mcp-use/mcp-use
https://github.com/nullplatform/meta-mcp-proxy
https://github.com/George5562/Switchboard
https://github.com/robertpelloni/mcp-tool-chainer
please examine all of those repos in depth and determine the best methods for each feature sot hat we have the best implementation
as there are many aggregators but i've tried most of them and none are ideal, each has some terrible flaw which makes it practically unusable and i always fall back to just having all mcps listed
and then the model never uses them anyway
And please also get inspiration for the dashboard design from the exmaple projects as well
The actual lesson for Borg
Borg should not be “an aggregator with a nicer search box.”

It should be a decision system.

The best implementation pattern per feature
1. Initial exposure
Best method:

expose only 5–6 permanent meta tools
no giant raw list
no giant suite tool unless forced by host limitations
Best starter set:

search_tools
load_tool
get_tool_schema
list_loaded_tools
unload_tool
run_code
That part we already documented — and I think it’s correct.

2. Search / discovery
Best method:

search must return ranked, compact, model-friendly results
each result should include:
why it matched
what task it is best for
typical latency/cost
whether it requires spawning a local binary
one short example
If search returns raw tool dumps, it’s dead on arrival.

3. Loading
Best method:

auto-load on high confidence
don’t force the model through a ceremonial “search then load then call” dance every time
In other words:

low confidence → show ranked options
high confidence → load silently and proceed
only ask for explicit selection when ambiguous
That is much better than making load_tool a mandatory step for everything.

4. Deferred binary startup
Best method:

index metadata without spawning binaries
spawn only on:
first execution
live schema hydration
explicit warm-up / health check
pinned always-warm profiles
This matters a lot because spawn latency trains the model not to bother.

5. Unloading
Best method:

unload by LRU idle policy
dehydrate schema before full eviction
keep loaded set tiny
I’d keep:

loaded metadata soft cap around 16
hard cap around 24
active binaries around 4
If loaded state grows much beyond that, you’re back to the same clutter problem in a different hat.

6. Tool grouping
Best method:

group for humans/operators
do not rely on grouping alone for model routing
Switchboard’s suite idea is useful, but suite→introspect→subtool is still extra model work.
Groups help organization; they do not solve tool selection by themselves.

7. Code mode
Best method:

use code mode as the escape hatch for complex multi-tool work
not as the universal first path
Lootbox is valuable here:

code mode is great once the right capability set is selected
it is not a replacement for discovery/routing
it is a multiplier after routing, not before it
8. Profiles
Best method:

let users pick or auto-detect task profiles
e.g.:
web research
repo coding
browser automation
local ops
database
Profiles are the only sane middle ground between:

“everything visible”
“search from scratch every time”
9. Observability
Best method:

show:
what the model searched for
what it loaded
what it ignored
load latency
spawn latency
failure reasons
eviction events
Without that, you can’t tell whether the model is failing because:

ranking is bad
descriptions are bad
latency is bad
prompt incentives are bad
My strongest conclusion after looking at these repos
If the goal is best practical implementation, the winning design is hybrid:

MCPM for install/import/profile/client-config ideas
Switchboard for lazy child spawning + tiered exposure ideas
Lazy MCP / Claude Lazy Loading for registry-first context minimization logic
Lootbox for code-mode execution once capabilities are chosen
Meta MCP Proxy for search/proxy minimalism
mcp-agent / mcp-use for lifecycle, pooling, auth, inspector, and operator-grade ergonomics
But the actual Borg implementation should be:

Borg’s best pattern
Tiny permanent meta-tool surface
Good ranked search
Silent auto-load when confidence is high
Deferred binary startup
Small active loaded set
Profiles for common workflows
Code mode for multi-step execution
Strong observability so the routing can be improved
That is very different from “just aggregate everything.”

The uncomfortable truth
A lot of aggregators are optimized for:

architecture elegance
token reduction on paper
protocol correctness
…but not for:

model behavior under uncertainty
That’s why they feel unusable in practice.

The model does not reward elegant indirection.
It rewards:

obvious affordances
low latency
high-confidence selection
small visible choice sets
What I’d do next
The next high-value move is not more repo comparison.

It’s to encode this as a Borg design rule:

The model should almost never face more than a handful of visible choices, and it should almost never be forced to manually perform the full discovery workflow when the system already knows the likely best capability.

If you want, I can do the next step right now:

update the research memo with a new section on “Why aggregators fail in practice”
and turn that into a concrete Borg implementation spec:
ranking model
confidence thresholds
auto-load rules
unload rules
profile behavior
UI/operator telemetry
That would turn the pain point into an implementable blueprint instead of just a correct complaint.
Done
I converted the “aggregators fail in practice” insight into durable repo guidance instead of leaving it floating in chat.

Updated files
MCP_ROUTER_REFERENCE_EVALUATION_2026-03-07.md

Added the second repo-set synthesis for:
mcpm.sh
Switchboard
claude-lazy-loading
lazy-mcp
lootbox
mcp-agent
mcp-use
meta-mcp-proxy
plus notes on the two 404 repos
Added a new section on why most aggregators fail in practice
Added a Borg implementation blueprint
Added a best method by feature matrix mapping each capability to the strongest upstream reference
CHANGELOG.md

Recorded the expanded router research memo and the new Borg hybrid blueprint guidance
What changed conceptually
The memo now makes the key distinction explicit:

most aggregators solve plumbing
Borg needs to solve selection friction
So the recommendation is now much sharper:

tiny permanent visible meta-tool set
ranked discovery, not just search
silent high-confidence auto-load
deferred binary startup
small active working set with eviction
code mode after routing
dashboard visibility into routing decisions, not just raw traffic
Tiny but important difference; that’s the gap between “technically aggregates” and “actually gets used.” Fancy plumbing is nice, but the model still needs a door handle.

Validation
I checked both edited files for issues.

MCP_ROUTER_REFERENCE_EVALUATION_2026-03-07.md — no errors
CHANGELOG.md — no errors
Todo status
 Synthesize repo findings
 Patch router research memo
 Update changelog if needed
Next likely move
The highest-value follow-up would be to translate this doc guidance into runtime behavior in core, especially:

reducing the effective loaded-tool cap
switching from FIFO-only behavior toward LRU/idle-aware eviction
adding confidence-based auto-load rules
exposing search/load/evict decisions in the dashboard inspector
That would turn the blueprint into working Borg behavior.
























Initial reference feature group categories:
MUST-HAVE BASIC FEATURES (START HERE):
mcp router/aggregator, mcp client which combines many mcp servers into one "master" mcp server, handles mcp session, instanciation, lifecycle auto start, restart, making sure it is available to multiple clients without multiple instances, making sure it does not time out and is responsive, measures latency, groups mcp servers into namespaces, groups tools into namespaces, enable/disable individual tools/servers, inspect mcp traffic, show activity, logging, mcp tool call chaining, TOON format, code mode, tool automatic renaming, tool reranking, tool semantic search, progressive tool disclosure, tool rag, context inspector, keep alive, heartbeat, mcp directory, automatic mcp install, environment variable and secrets management, mcp client config file autoconfig, auto write config file format, import/export configs, autodetect all configs on system, maintain/manage all configs on system, save sets of configs, wipe all client configs, set all client configs to x, etc, manage databases, handle auth. 
mcp server chains other mcps. 
can watch mcp traffic. 
reduce context language. 
reduce context syntax. 
resource management, lazy load, health checks, auto start, auto restart, timeout, deferred loading of binary until needed, latency, responsiveness, single-instance multiple clients, monitor memory usage, cpu usage, long running process. 
proxy for stdio as remote. 
proxy for remote as stdio. 
proxy for sse as streaming-http. 
proxy for streaming-http as stdio. 
manage oauth, bearer token. 
manage secrets, environment variables, paths. 
manage env variable expansion, .env files. 
mcp groups. 
tool groups. disable tools.
tool renaming, minimize context. 
tool semantic search/tool rag. 
tool lazy load, dynamic progressive tool disclosure. 
tools as code, code mode. 
mcp tool chaining. 
TOON format, context saving. 
inspect context, context makeup,
specify extremely high value tools to be disclosed always
loading/unloading after tool count
auto-load on high confidence




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

subagents, long running, subcontexts, timeouts, multiple models, mutiple cli, orchestration, harness/cli/tiu as mcp. 
multi-model debate, multi-model consensus, share context between models, multiple models pair programming with each other, architect-implementer. 

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


Outstanding! Magnificent! Insanely Great!!! Please continue to proceed as you recommend based on your ongoing analysis and assessments! Keep going! 