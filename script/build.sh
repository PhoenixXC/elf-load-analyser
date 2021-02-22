#!/bin/bash
readonly TARGET="ela"
readonly OUT_SUFFIX="k"
readonly WORK_DIR=$(pwd)
readonly GO_WORK="${WORK_DIR}/../"
readonly SRC_DIR="${WORK_DIR}/../pkg/modules/module/src"
echo "Bcc source dir: ${SRC_DIR}"

cd "${SRC_DIR}" || exit
echo "Merge files"
for line in *.cpp
do
  grep -E '^#include[[:blank:]]+"[[:print:]]+"$' "${line}" | awk '{print $2}' | tr -d '"' | awk '{
    if (NF > 0) {
      if ($1!="_dev.h") print $1
    }}' | xargs cat > "${line}"."${OUT_SUFFIX}"
  echo "" >> "${line}"."${OUT_SUFFIX}"
  sed '/^[[:blank:]]*#include[[:blank:]]\+"[[:print:]]\+"$/d' "${line}" | cat >> "${line}"."${OUT_SUFFIX}"
done

cd "${GO_WORK}" || exit
echo "Build binary..."
mkdir -p "target"
go build -o target/"${TARGET}" github.com/phoenixxc/elf-load-analyser/cmd
echo "Build ok, now you can use '$(pwd)/target/${TARGET}' run program"

cd "${SRC_DIR}" || exit
echo "Clean tmp files"
rm -- *."${OUT_SUFFIX}"