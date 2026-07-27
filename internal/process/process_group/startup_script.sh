#!/usr/bin/env bash

set -Eeuo pipefail

export -n _CHDIR_ OVERLAY_ROOT

if [[ -z "$_CHDIR_" ]]; then
	echo "缺少必须的环境变量 _CHDIR_" >&2
	exit 1
fi

if [[ $# -eq 0 ]]; then
	echo "缺少要执行的命令 (或没有加入 --)" >&2
	exit 1
fi

echo "程序: $*" >&2
echo "工作目录: $_CHDIR_" >&2
echo "Overlay根目录: $OVERLAY_ROOT" >&2

if [[ $$ -ne 1 ]]; then
	echo "启动脚本必须在 PID 1 下运行" >&2
	exit 1
fi

mount --make-rprivate /

# 如果要求挂载根目录的 overlay 文件系统，则执行挂载操作
# 将所有写入操作重定向到 OVERLAY_ROOT 目录，并将真正的根目录作为 lowerdir
if [[ -n "$OVERLAY_ROOT" ]]; then
	mkdir -p "${OVERLAY_ROOT}" "${OVERLAY_ROOT}.work" "${OVERLAY_ROOT}.lower" "${OVERLAY_ROOT}.merged"
	mount --bind / "${OVERLAY_ROOT}.lower"

	mount -t overlay overlay \
		-o "lowerdir=${OVERLAY_ROOT}.lower,upperdir=$OVERLAY_ROOT,workdir=${OVERLAY_ROOT}.work" \
		"${OVERLAY_ROOT}.merged"

	cd "${OVERLAY_ROOT}.merged"
	pivot_root . tmp

	mount -t proc proc /proc
	mount --move /tmp/dev /dev
	mount --move /tmp/sys /sys
	mount --move /tmp/run /run
	umount -l /tmp
else
	mount -t proc proc /proc
fi
mount -t tmpfs tmpfs /tmp

cd "$_CHDIR_"

exec "$@"
