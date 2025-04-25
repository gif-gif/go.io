#!/bin/bash

# 脚本描述
echo "请选择要更新的版本号类型："
echo "1. 更新主版本号（如 v1.0.0 -> v2.0.0）"
echo "2. 更新测试版本号（如 v0.1.0 -> v0.2.0）"
echo "3. 更新修订版本号（如 v0.0.1 -> v0.0.2）"
echo "默认使用修订版本号，直接按回车执行。"

# 读取用户输入
read -p "请输入选项（1/2/3，默认 3）：" choice

# 设置默认值为 3（修订版本号）
choice=${choice:-3}

# 获取当前最新的 Tag
latest_tag=$(git describe --tags --abbrev=0 2>/dev/null)

echo "当前仓库最新tag版本号为：" $latest_tag

# 如果没有找到 Tag，默认从 v0.0.0 开始
if [ -z "$latest_tag" ]; then
  latest_tag="v0.0.0"
fi

# 提取版本号（去掉开头的 'v'）
version=${latest_tag#v}

# 将版本号拆分为主版本号、测试版本号和修订版本号
IFS='.' read -r major minor patch <<< "$version"

# 根据用户选择更新版本号
case $choice in
  1)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  2)
    minor=$((minor + 1))
    patch=0
    ;;
  3)
    patch=$((patch + 1))
    ;;
  *)
    echo "无效选项，默认使用修订版本号。"
    patch=$((patch + 1))
    ;;
esac

# 拼接新的版本号
new_tag="v${major}.${minor}.${patch}"

echo "本次即将提交的tag版本号为:" $new_tag

## 打 Tag
git tag "$new_tag"
git push origin "$new_tag"

# 输出最新的 Tag
echo "新的 Tag 已创建：$new_tag"