<a name="unreleased"></a>
## [Unreleased]


<a name="v0.9.0"></a>
## [v0.9.0] - 2025-11-26
### Chore
- 更新 Makefile，添加 changelog 相关命令并初始化配置
- tui messages type package definition
- Add change log template and configuration file

### Fix
- 修正代办列表帮助信息，一些emoji会导致界面渲染错位

### Refactor
- 重构 TodoList 页面，使用 bubbles 组件库
- 移除应用初始化逻辑，直接在运行时加载任务


<a name="v0.8.0"></a>
## [v0.8.0] - 2025-10-29
### Refactor
- Switched TUI framework from tview to bubbletea.


<a name="v0.7.3"></a>
## [v0.7.3] - 2025-10-29
### Chore
- update installation path to ~/.local/bin and log config file path correctly
- update log filename path to include logs directory

### Refactor
- set default tasks save path and format JSON output for task items


<a name="v0.7.0"></a>
## [v0.7.0] - 2025-10-29

<a name="v0.7.2"></a>
## [v0.7.2] - 2025-10-29
### Chore
- update installation path to ~/.local/bin and log config file path correctly


<a name="v0.7.1"></a>
## [v0.7.1] - 2025-10-29
### Chore
- update log filename path to include logs directory
- Update config package and add default configuration file

### Feat
- add Makefile and implement versioning and system info logging

### Refactor
- Update config path for xdg compliance and init xdg paths.


<a name="v0.6.2"></a>
## [v0.6.2] - 2024-07-03
### Chore
- update license to MIT License


<a name="v0.6.1"></a>
## [v0.6.1] - 2024-07-02
### Refactor
- improve file saving and logging


<a name="v0.6.0"></a>
## [v0.6.0] - 2024-06-28
### Feat
- save tasks to file


<a name="v0.5.0"></a>
## v0.5.0 - 2024-06-28
### Add
- git ignore

### Chore
- add log、config pkg
- add TestSwitchPagesAndContent

### Feat
- add todo-list page
- add config
- page and content switch test
- add logger
- app base

### Init
- kong tools init

### Refactor
- limit log backup

### Style
- update welcome page


[Unreleased]: https://github.com/yizhixiaokong/kongtools/compare/v0.9.0...HEAD
[v0.9.0]: https://github.com/yizhixiaokong/kongtools/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/yizhixiaokong/kongtools/compare/v0.7.3...v0.8.0
[v0.7.3]: https://github.com/yizhixiaokong/kongtools/compare/v0.7.0...v0.7.3
[v0.7.0]: https://github.com/yizhixiaokong/kongtools/compare/v0.7.2...v0.7.0
[v0.7.2]: https://github.com/yizhixiaokong/kongtools/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/yizhixiaokong/kongtools/compare/v0.6.2...v0.7.1
[v0.6.2]: https://github.com/yizhixiaokong/kongtools/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/yizhixiaokong/kongtools/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/yizhixiaokong/kongtools/compare/v0.5.0...v0.6.0
