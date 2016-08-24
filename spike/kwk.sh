#!/bin/bash

CMD="$1"
INPUT="$2"
INPUT2="$3"
#API_ROOT="https://kwk.co/api/v1/"
#API_ROOT="http://localhost:8080/api/v1/"
API_ROOT="http://kwk.loc/api/v1/"
#API_ROOT="http://192.168.99.100/api/v1/"
CACHE=/Users/$USER/Library/Caches/kwk.cache
PREFERENCES=/Users/$USER/Preferences/kwk.preferences

orange=`tput setaf 9`
white=`tput setaf 15`
error=`tput setaf 1`
blue=`tput setaf 33`
subdued=`tput setaf 7`
reset=`tput sgr0`
bg=`tput setab 33`
prompt="\n${blue}kwk>${reset}"

function addToCache {
    local KEY=$1
    local URL=$2
    if [ ! -d ${CACHE} ]; then
        mkdir ${CACHE}
    fi
    echo ${URL} > ${CACHE}/${KEY}
}

function kwkHelp {
        printf "${bg}${white}<==============> WELCOME TO KWK! <==============>${reset}${blue}\n\n"
        printf "Commands:\n\n"
        printf "  kwk <KEY,..>                    - Open the KEY(S) and navigate.\n" \
        " If the key is a priority key then open the first link, or if its a cluster key then"
        " open all the links"
        printf "  kwk open tag <KEY> [PAGE]       - Open all links with a given tag in desc in pages of 5"
        printf "  kwk create <URI> [DESIRED KEY]  - Create a KEY and customize.\n"
        printf "  kwk suggest <URI>               - Returns a list of suggested KEYS.\n"
        printf "  kwk tag <KEY> [TAG1,TAG2,..]    - Add tags to a key.\n"
        printf "  kwk untag <KEY> [TAG1,TAG2,..]  - Remove tags from a key.\n"
        printf "  kwk lock <KEY> [PIN]            - Lock a key with a pin.\n"
        printf "  kwk update <OLD KEY> <NEW KEY>  - Made a mistake? Update a key.\n"
        printf "  kwk detail <KEY>                - Get KEY details and info.\n"
        printf "  kwk covert <KEY>                - Open the KEY in covert (incognito mode).\n"
        printf "  kwk get <KEY>                   - Get the URI without navigating and copies it.\n"
        printf "  kwk search [today][yesterday]   - Search remote KEYS.\n"

        printf "\nAnalytics:\n\n"
        printf "  kwk stats [today, yesterday]    - Get stats summary for all keys\n"
        printf "  kwk stats <KEY>                 - Get stats for a key plus sub keys.\n"

        printf "\nAccount:\n\n"
        printf "  kwk login <SECRET> [dontsave]   - Login with secret key.\n"
        printf "  kwk logout                      - Clears locally cached secret key.\n"
        printf "  kwk subscribe <KEY>             - Reserve your own personalised top-level key.\n"

        printf "\nAdmin/Config:\n\n"
        printf "  kwk cache ls                    - List locally cached keys.\n"
        printf "  kwk cache clear                 - Clears any locally cached data.\n"
        printf "  kwk upgrade                     - Downloads and upgrades kwk cli client.\n"
        printf "  kwk config warn  [on|off]       - Warns if attempting to open dodgy key.\n"
        printf "  kwk config quiet [on|off]       - Prevents links from being printed to console.\n"
        printf "  kwk version                     - Enough said.\n"

        printf "\n<===============================================>\n\n"
        printf "${reset}"
}

function encode {
    local url=$1
    local desiredKey=$2
    local body=""

    if [[ ! -z ${desiredKey} ]]; then
       body='{"url":"'${url}'","key":"'${desiredKey}'"}'
    else
        body='{"url":"'${url}'"}'
    fi

    local response=$(curl -s -H "Content-Type: application/json" -H "x-kwk-key: 59364212-aeb2-4100-bf0e-c418ef230529" -XPOST -d ${body} "${API_ROOT}hash")
    local key=$(echo ${response} | jq .key | tr -d '"')
    local err=$(echo ${response} | jq .error | tr -d '"')

    if [[ "${err}" == "null" ]]; then
        printf "\n\n Your kwk code: ${bg}${white}  ${key}  ${reset}"
        echo ${key} | tr -d '\n' | pbcopy
        printf "\n\n${subdued}(copied)${reset}\n\n"
        addToCache ${key} ${url}
    else
         printf "${prompt} ${error}${err}${reset}\n"
         exit 1
    fi
}

function decode {
    local input=$1
    local url=""
    if [[ -f ${CACHE} ]]; then
       url=$(cat ${CACHE})
       nohup curl -s "${API_ROOT}hash/${input}" >/dev/null 2>&1 &
    else
       url=$(curl -s "${API_ROOT}hash/${input}" | jq .url | tr -d '"')
       if [[ "${url}" == "null" ]]; then
            printf "${prompt} kwk code doesn't exist\n"
            exit 1
       fi
       addToCache ${input} ${url}
    fi
    echo "${url}"
}

function upgrade {
    local distributionUri=$(decode kwk.dist)
    local installPath=$(decode mac.kwk)
    #download in future
    sudo cp "${distributionUri}" "${installPath}"
    printf "${prompt} Upgraded!" | tr -d '\n'
    printf "\n"
}

function openPath {
    local command=$1
    local path=$2
    printf "${path}"
     if [[ "${command}" == "get" ]]; then
      echo ${path} | tr -d '\n' | pbcopy
    elif [[ "${command}" == "covert" ]]; then
        /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --incognito ${path}
    else
       open ${path}
    fi
}

case "$CMD" in
"create" | "s")
    printf "${prompt} Saving ${INPUT}" | tr -d "\n"
    encode ${INPUT} ${INPUT2}
    ;;
"update")
        printf "${blue}Coming soon${reset}\n"
    ;;
"upgrade")
    upgrade
    ;;
"version")
    printf "${prompt} 0.0.1" | tr -d "\n"
    printf "\n"
    ;;
"cache")
  case "$INPUT" in
        "list" | "ls")
            printf "${prompt}${blue} Your Cache in ${CACHE}" | tr -d "\n"
            printf "\n"
            ls ${CACHE}
            printf "${reset}\n"
        ;;
        "clear")
            rm -fr "${CACHE}"
            printf "${prompt} Cache cleared." | tr -d "\n"
            printf "\n"
        ;;
        *)
          kwkHelp
        ;;
    esac
  ;;
*)
    if [[ -z "${CMD}" ]]; then
         kwkHelp
         exit 0
    fi

    if [[ -z ${INPUT} ]]; then
        INPUT=${CMD}
        CMD="open"
    fi

    DECODED=$(decode "${INPUT}")
    if [[ $? == 0 ]]; then
        openPath ${CMD} ${DECODED}
    else
        printf "${prompt} That command or key doesn't exit"
    fi
    printf "\n"
    ;;
esac