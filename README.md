### cmt 
> git 提交说明规范化工具，类似于 Commitizen 这个工具


#### 0. 由来
* 安装 commitizen 工具的时候发现其需要安装的依赖太多
* 工具的主要功能是提交说明的格式化及校验
* 这个功能用golang来实现比较简单
* `github.com/AlecAivazis/survey`库可以很好的支持这个工具需求的实现
---
#### 1. 输入信息
* 需要输入此次提交的类型，`type`，目前的类型（`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `revert`, `WIP`）
  * feat: A new feature 
  * fix: A bug fix, 
  * docs: Documentation only changes, 
  * style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc), 
  * refactor: A code change that neither fixes a bug nor adds a feature, 
  * perf: A code change that improves performance, 
  * test: Adding missing tests, 
  * chore: Changes to the build process or auxiliary tools and libraries such as documentation generation, 
  * revert: Revert to a commit, 
  * WIP: Work in progress,
* 此次提交的代码影响的模块范围，`scope`
  * 输入会有校验，由小写字母和数字组成，2-20个字符
* 提交的说明信息，`message`
* 项目的版本信息，`version`
* 项目对应的 teambition 单子id信息，`tb_num`
---
#### 2. 功能
* 根据输入的各项信息，组成格式化的提交信息
  * 格式类似：`feat(core): test commit message(v1.3.0-BF-11543)`
* 根据上一次的提交来设置每一项的默认值，减少每次有些和上次提交重复的信息
  * 通过调用git命令：`git log --pretty=format:"%s" --no-merges -1`
  * 根据返回来提取每项关键信息，然后填入默认值
* 提交说明输入完成后会直接调用`git commit`命令进行提交
---
#### 3. 使用
* 按照需求进行改动每一项输入的逻辑改动
* 构建可执行文件
  * 使用go命令构建：`go build -o git-cm .`
  * 支持`ldflags`：`go build -ldflags="-X main.zh=t" -o git-cm .`, `type`的说明支持中文
  * 如果想要使用 `git` 的 `subcommands`, 类似：`git mytool`，需要将构建的二进制名称进行重命名成`git-mytool`
* 将可执行文件加入到`PATH`中（zsh示例）
  * echo "{cmtpath}:$PATH" > ~/.zshrc
  * source ~/.zshrc
---
#### 4. feature
* ~~增加type配置~~
* ~~增加提交说明格式配置~~
* ~~增加输入项配置~~
---
#### 5. Enjoy it!
* 欢迎fork，代码逻辑较为简单，可以根据自身需要进行定制化开发