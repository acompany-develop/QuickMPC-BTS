#!/bin/bash
# ----- NOTE: 以下は変更不要 ----- #

# エラーが発生しても処理を続けるモードに切り替え
set +e

# BuildKitを使用して高速にBuild
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

status=0

# 引数が正しくない時に発火する関数
usage_exit() {
    echo "Usage: $1 'build' or 'run' or ''" 1>&2
    exit 1
}

# build()の定義を強要するために未定義関数として宣言
build() {
    echo "build()が定義されていません"
    status=$(($status + 1))
}

# run()の定義を強要するために未定義関数として宣言
run() {
    echo "run()が定義されていません"
    status=$(($status + 1))
}

# setup()の定義を強要するために未定義関数として宣言
setup() {
    echo "setup()が定義されていません"
    status=$(($status + 1))
}

# teardown()の定義を強要するために未定義関数として宣言
teardown() {
    echo "teardown()が定義されていません"
    status=$(($status + 1))
}


if [ $# -eq 1 ]; then
    source $1
    status=$(($status + $?))
    if [ $status -gt 0 ]; then
        echo $status >>result
    else
        setup
        status=$(($status + $?))
        build
        status=$(($status + $?))
        run
        status=$(($status + $?))
        # NOTE: `docker-compose up`はCMDやENTRYPOINTで異常終了してもexitステータスが`0`になってしまうので別途exitステータスを集積する
        run_status=$(docker-compose -f docker-compose.yml ps -q | tr -d '[ ]' | xargs docker inspect -f '{{ .State.ExitCode }}' | grep -v 0 | wc -l | tr -d '[ ]')
        status=$(($status + $run_status))
        teardown
        status=$(($status + $?))
        echo $status >>result
    fi
else
    source $1
    status=$(($status + $?))

    if [ $status -gt 0 ]; then
        echo $status >>result
    else
        case $2 in
        build)
            build
            status=$(($status + $?))
            echo $status >>result
            ;;
        run)
            setup
            status=$(($status + $?))
            run
            status=$(($status + $?))
            # NOTE: `docker-compose up`はCMDやENTRYPOINTで異常終了してもexitステータスが`0`になってしまうので別途exitステータスを集積する
            run_status=$(docker-compose -f docker-compose.yml ps -q | tr -d '[ ]' | xargs docker inspect -f '{{ .State.ExitCode }}' | grep -v 0 | wc -l | tr -d '[ ]')
            status=$(($status + $run_status))
            teardown
            status=$(($status + $?))
            echo $status >>result
            ;;
        *)
            usage_exit
            ;;
        esac
    fi
fi

# DEBUG: 最後にstatusを表示
echo "["$1"] status:" $status

# エラーが発生したらすぐに終了するモードに戻す
set -e

exit 0
