echo "enter version: "

read version

options=("build api" "build web")

select opt in "${options[@]}"
do
    case $opt in
        "build api")
            # remove docker image if exists
            docker rmi -f chatgpt-plus-api:$version
            # build docker image for chatgpt-plus-go
            docker build -t chatgpt-plus-api:$version -f dockerfile-api-go ../
            ;;
        "build web")
            docker rmi -f chatgpt-plus-web:$version
            docker build -t chatgpt-plus-web:$version -f dockerfile-vue ../
            ;;
        *) echo "invalid option $REPLY";;
    esac
    break
done

