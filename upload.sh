#!/bin/bash
# 自动提交并上传到 GitHub 的脚本
# 用法示例： ./upload.sh "上传第三节课作业"

if [ -z "$1" ]; then
  echo "⚠️  请输入提交说明，例如：./upload.sh '更新第四节课'"
  exit 1
fi

# 检查是否有未完成的 rebase
if [ -d ".git/rebase-merge" ] || [ -d ".git/rebase-apply" ]; then
  echo "⚠️ 检测到未完成的 rebase，请先运行：git rebase --continue 或 git rebase --abort"
  exit 1
fi

# 拉取最新远程分支
git pull --no-rebase

# 添加、提交并推送更改
git add .
git commit -m "$1"
git push

echo "✅ 上传成功！GitHub 仓库已更新 🚀"
