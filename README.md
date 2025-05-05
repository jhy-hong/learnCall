## git相关操作

### 标签的使用

####  一、标签的基本操作

1. 创建标签​：

```bash
# 轻量标签（无备注）

git tag <tag-name>

# 附注标签（带备注）

git tag -a <tag-name> -m "tag message"

# 示例

git tag v1.0.0
git tag -a v1.0.1 -m "Release version 1.0.1"
```

2. 查看标签​：
   ```bash
   # 列出所有标签
   
   git tag
   
   # 查看标签详情
   
   git show <tag-name>
   ```

3. 推送标签到远程仓库​：
   ```bash
   # 推送单个标签
   
   git push origin <tag-name>
   
   # 推送所有本地标签
   
   git push origin --tags
   ```

3. 删除标签​：
   ```bash
   # 删除本地标签
   
   git tag -d <tag-name>
   
   # 删除远程标签
   
   git push origin :refs/tags/<tag-name>
   ```

#### 二、提交代码到标签的后续操作
   由于标签本身不可变，若要将新代码关联到某个标签，通常有以下两种方式：

1. 覆盖（强制更新）标签​
   如果标签尚未推送到远程仓库，或者你确定需要修改标签指向的提交：

```bash
# 1. 删除旧标签（本地）

git tag -d <tag-name>

# 2. 创建新标签指向当前提交

git tag <tag-name>

# 3. 强制推送到远程

git push origin :refs/tags/<tag-name>  # 删除远程旧标签
git push origin <tag-name>            # 推送新标签
```

2. 基于标签创建新分支​
   如果希望基于某个标签继续开发：

```bash
# 1. 基于标签创建新分支

git checkout -b <new-branch> <tag-name>

# 2. 在新分支上提交代码

git add .
git commit -m "new changes"

# 3. 推送分支到远程

git push origin <new-branch>

# 4. 可选：为新提交打新标签（如 v1.0.1-hotfix）

git tag -a v1.0.1-hotfix -m "Hotfix for v1.0.1"
git push origin v1.0.1-hotfix
```



#### 三、最佳实践

1. **标签不可变原则**：
   - 标签应指向不可变的版本（如正式发布版本），避免频繁修改。
   - 如果需要修复旧版本的 Bug，建议：
     1. 基于旧版本标签创建分支（如 `git checkout -b hotfix-v1.0.1 v1.0.1`）。
     2. 在分支上修复并提交。
     3. 为新修复版本打新标签（如 `v1.0.2`）。
2. **语义化版本命名**：
   - 使用 SemVer 规范（如 `v1.2.3`），便于版本管理。
3. **强制更新标签的风险**：
   - 若远程标签已被他人使用，强制更新可能导致协作问题。建议通过新版本解决。

#### 四、示例场景

假设当前有一个标签 `v1.0.0`，你需要修复该版本的 Bug：

```bash
# 1. 基于标签创建分支

git checkout -b hotfix-v1.0.0 v1.0.0

# 2. 修复 Bug 并提交

git add .
git commit -m "Fix critical bug"

# 3. 打新标签

git tag v1.0.1

# 4. 推送分支和标签

git push origin hotfix-v1.0.0
git push origin v1.0.1
```



### 分支的相关使用与操作

Git 的分支（Branch）是代码开发中非常重要的功能，允许你在独立的环境中开发新功能或修复 Bug，而不会影响主分支（如 `master` 或 `main`）。以下是分支的常用操作及流程：

### 一、分支的基本操作

#### 1. **查看分支**

```bash
# 查看本地分支（当前分支前会标 *）
git branch

# 查看远程分支
git branch -r

# 查看所有分支（本地 + 远程）
git branch -a
```

#### 2. **创建分支**

```bash
# 创建新分支（但不会自动切换）
git branch <branch-name>

# 创建并切换到新分支（推荐）
git checkout -b <branch-name>
```

#### 3. **切换分支**

```bash
# 切换到已有分支
git checkout <branch-name>

# 切换到远程分支（会自动创建本地分支跟踪远程分支）
git checkout -b <local-branch-name> origin/<remote-branch-name>
```

#### 4. **删除分支**

```bash
# 删除本地分支（需先切换到其他分支）
git branch -d <branch-name>

# 强制删除未合并的分支
git branch -D <branch-name>

# 删除远程分支
git push origin --delete <remote-branch-name>
```

#### 5. **推送分支到远程仓库**

```bash
# 推送本地分支到远程（并建立关联）
git push -u origin <branch-name>

# 后续推送可简写为
git push
```

### 二、合并分支（Merge）

#### 1. **合并分支到当前分支**

```bash
# 1. 切换到目标分支（如主分支）
git checkout main

# 2. 合并其他分支（如 feature/login）
git merge feature/login
```

#### 2. **合并冲突解决**

如果合并时发生冲突，需手动解决：

```bash
# 1. 冲突文件会标出冲突位置（<<<<<<< 和 >>>>>>>）
# 2. 编辑文件，保留需要的代码并删除冲突标记
# 3. 标记冲突已解决
git add <file-name>

# 4. 提交合并结果
git commit -m "Merge branch 'feature/login' into main"
```

#### 3. **合并模式**

- **快进合并（Fast-Forward）**：如果目标分支没有新提交，直接移动指针（无合并提交）。
- **三方合并（Three-Way Merge）**：目标分支有新提交，生成一个新的合并提交。

### 三、变基分支（Rebase）

变基（Rebase）可以将分支的提交历史“整理”成一条直线，使历史更清晰（但需谨慎使用）。

#### 1. **变基操作**

```bash
# 1. 切换到待整理的分支（如 feature/login）
git checkout feature/login

# 2. 将当前分支变基到目标分支（如 main）
git rebase main

# 3. 解决可能的冲突（类似 merge）
# 4. 切换回主分支并合并
git checkout main
git merge feature/login
```

#### 2.**注意**

- 变基会修改提交历史，**不要对已推送到远程的分支执行变基**（除非确认只有自己使用该分支）。

### 四、分支管理最佳实践

#### 1. **分支命名规范**

- 功能分支：`feature/<功能名>`（如 `feature/login`）
- Bug 修复分支：`bugfix/<问题描述>`（如 `bugfix/login-error`）
- 发布分支：`release/<版本号>`（如 `release/v1.0.0`）

#### 2. **分支生命周期**

- 开发新功能时，从主分支创建新分支，开发完成后合并回主分支并删除。
- 长期分支（如 `dev`）可用于集成测试。

#### 3. **保持分支同步**

```bash
# 从远程拉取最新代码并合并到当前分支
git pull origin main

# 或者先拉取再合并
git fetch origin
git merge origin/main
```

#### 4. **频繁提交与合并**

- 小步提交，避免一次提交过多改动。
- 定期将主分支的更新合并到开发分支，减少冲突概率。

### 五、常见工作流示例

#### 1. **开发新功能**

```bash
# 1. 创建并切换到新分支
git checkout -b feature/add-search

# 2. 开发代码并提交
git add .
git commit -m "Add search functionality"

# 3. 推送分支到远程
git push -u origin feature/add-search

# 4. 合并到主分支（通过 Pull Request/Merge Request）
git checkout main
git merge feature/add-search

# 5. 删除已合并的分支
git branch -d feature/add-search
```

#### 2. **修复线上 Bug**

```bash
# 1. 从主分支创建修复分支
git checkout -b hotfix/login-error

# 2. 修复问题并提交
git add .
git commit -m "Fix login timeout error"

# 3. 合并到主分支并打标签
git checkout main
git merge hotfix/login-error
git tag v1.0.1
git push origin v1.0.1
```

假设你误将代码提交到了标签 `v1.0.0` 对应的提交，但主分支未同步：

```bash
# 1. 查看标签对应的提交ID
git show-ref --tags
# 输出：abcdef12 refs/tags/v1.0.0

# 2. 切换到主分支并重置
git checkout master
git reset --hard abcdef12

# 3. 强制推送主分支
git push --force origin master
```

