#!/bin/sh

gitClone(){
  [ $VERBOSE = 0 ] || printf "Cloning repository...\n"
  ssh-keyscan github.com >> ~/.ssh/known_hosts 2> /dev/null
  echo "$(cat ~/.ssh/known_hosts | sort | uniq)" > ~/.ssh/known_hosts
  FLAG='&> /dev/null'
  [ $VERBOSE = 0 ] || unset -v FLAG
  git clone git@github.com:qdm12/updated.git . $FLAG
  unset -v FLAG
}

gitPull(){
  [ $VERBOSE = 0 ] || printf "Pulling files...\n"
  ssh-keyscan github.com >> ~/.ssh/known_hosts 2> /dev/null
  echo "$(cat ~/.ssh/known_hosts | sort | uniq)" > ~/.ssh/known_hosts
  FLAG='&> /dev/null'
  [ $VERBOSE = 0 ] || unset -v FLAG
  git pull $FLAG
  unset -v FLAG
}

gitPush(){
  [ $VERBOSE = 0 ] || printf "Pushing files...\n"
  ssh-keyscan github.com >> ~/.ssh/known_hosts 2> /dev/null
  echo "$(cat ~/.ssh/known_hosts | sort | uniq)" > ~/.ssh/known_hosts
  FLAG='&> /dev/null'
  [ $VERBOSE = 0 ] || unset -v FLAG
  git add .
  git commit -m "Update of $(date)" $FLAG
  git push $FLAG
  unset -v FLAG
}

buildMaliciousHostnames(){
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Building NSA hostnames\n"
  hostnames=$(wget -qO- https://raw.githubusercontent.com/dyne/domain-list/master/data/nsa | \
  sed '/\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)/d' | \
  sed 's/\(\r\)//g')$'\n'$( \
  wget -qO- "https://raw.githubusercontent.com/Cauchon/NSABlocklist-pi-hole-edition/master/HOSTS%20(including%20excessive%20GOV%20URLs)" | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^127.0.0.1\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(0.0.0.0 \)//g')$'\n'$( \
  wget -qO- https://raw.githubusercontent.com/CHEF-KOCH/NSABlocklist/master/HOSTS | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^127.0.0.1\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(0.0.0.0 \)//g')
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  hostnames=$(echo "$hostnames" | tr '[:upper:]' '[:lower:]' | sort | uniq)
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $(($COUNT_BEFORE-$COUNT_AFTER)) duplicates from $COUNT_BEFORE hostnames\n"
  echo "$hostnames" > files/nsa-hostnames.updated
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v hostnames
  unset -v COUNT_BEFORE
  unset -v COUNT_AFTER
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Building hostnames\n"
  hostnames=$(wget -qO- https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^127.0.0.1\)\|\(^255.255.255.255\)\|\(^::1\)\|\(^fe80\)\|\(^ff00\)\|\(^ff02\)\|\(^0.0.0.0 0.0.0.0\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(0.0.0.0 \)//g')$'\n'$( \
  wget -qO- https://raw.githubusercontent.com/k0nsl/unbound-blocklist/master/blocks.conf | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^local-data\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(local-zone: \"\)\|\(\" redirect\)//g')$'\n'$( \
  wget -qO- https://raw.githubusercontent.com/notracking/hosts-blocklists/master/domains.txt | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(::$\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(address=\/\)\|\(\/0.0.0.0$\)//g')$'\n'$( \
  wget -qO- https://raw.githubusercontent.com/notracking/hosts-blocklists/master/hostnames.txt | \
  sed '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^::\)/d' | \
  sed 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(^0.0.0.0 \)//g')
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  hostnames=$(echo "$hostnames" | tr '[:upper:]' '[:lower:]' | sort | uniq | sed '/\(psma01.com.\)\|\(psma02.com.\)\|\(psma03.com.\)/d')
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $(($COUNT_BEFORE-$COUNT_AFTER)) duplicates from $COUNT_BEFORE hostnames\n"
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  hostnames=$(echo "$hostnames" | tr '[:upper:]' '[:lower:]' | sed '/\(maxmind.com\)\|\(ipinfo.io\)/Id')
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$hostnames" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $(($COUNT_BEFORE-$COUNT_AFTER)) allowed entries\n"
  echo "$hostnames" > files/malicious-hostnames.updated
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v hostnames
  unset -v COUNT_BEFORE
  unset -v COUNT_AFTER
}

buildMaliciousIPs(){
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Building malicious IPs\n"
  ips=$(wget -qO- https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt | grep -v "#" | grep -v -E "\s[1-2]$" | cut -f 1)$'\n'$( \
  wget -qO- https://iplists.firehol.org/files/firehol_level1.netset | grep -v "#" | grep -v -E "\s[1-2]$" | cut -f 1)
  ips=$(echo "$ips" | sed '/\(^127.\)\|\(^10.\)\|\(^172.1[6-9].\)\|\(^172.2[0-9].\)\|\(^172.3[0-1].\)\|\(^192.168.\)/d')
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  ips=$(echo "$ips" | sort | uniq)
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $(($COUNT_BEFORE-$COUNT_AFTER)) duplicates from $COUNT_BEFORE ips\n"
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  ips=$(echo "$ips")
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $(($COUNT_BEFORE-$COUNT_AFTER)) allowed entries\n"
  echo "$ips" > files/malicious-ips.updated
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v ips
  unset -v COUNT_BEFORE
  unset -v COUNT_AFTER
}

extendMaliciousIPs(){
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Extending malicious IPs\n"
  ips=$(cat files/malicious-ips.updated)
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  while read hostname; do
    [ $VERBOSE = 0 ] || printf "Resolving $hostname..."
    resolvedips=$(dig +short +time=1 +tries=1 "$hostname" | grep -v "no servers could be reached" | grep -oE "\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b" | grep -vE "(^127\.)|(^10\.)|(^172\.1[6-9]\.)|(^172\.2[0-9]\.)|(^172\.3[0-1]\.)|(^192\.168\.)")
    for ip in $resolvedips; do
      ips=$'\n'$ip
      printf " $ip"
    done
    [ $VERBOSE = 0 ] || printf "\n"
  done <files/malicious-hostnames.updated
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Added $((COUNT_AFTER-$COUNT_BEFORE)) IPs from resolved malicious hostnames\n"
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  ips=$(echo "$ips" | sort | uniq)
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $((COUNT_BEFORE-$COUNT_AFTER)) duplicates from $COUNT_BEFORE ips\n"
  [ $VERBOSE = 0 ] || COUNT_BEFORE=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  ips=$(echo "$ips")
  [ $VERBOSE = 0 ] || COUNT_AFTER=$(echo "$ips" | sed '/^\s*$/d' | wc -l)
  [ $VERBOSE = 0 ] || printf "Removed $((COUNT_BEFORE-$COUNT_AFTER)) allowed entries\n"
  echo "$ips" > files/malicious-ips.updated
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v ips
  unset -v COUNT_BEFORE
  unset -v COUNT_AFTER
}

buildNamedRoot(){
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Building DNS Named root\n"
  MD5SUM=$(wget -qO- https://www.internic.net/domain/named.root.md5)
  wget -q https://www.internic.net/domain/named.root -O files/named.root.temp
  if [ "$(md5sum files/named.root.temp | cut -d " " -f 1)" != "$MD5SUM" ]; then
    printf "named.root MD5 checksum mismatch !\n"
    echo "1" > health
    rm files/named.root.temp
    return
  fi
  rm files/named.root.updated
  mv files/named.root.temp files/named.root.updated
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v MD5SUM
}

buildRootAnchors(){
  [ $VERBOSE = 0 ] || T_START=$(date +%s)
  [ $VERBOSE = 0 ] || printf "Building DNS Root anchors\n"
  wget -q https://data.iana.org/root-anchors/root-anchors.xml -O files/root-anchors.xml.temp
  if [ "$(sha256sum files/root-anchors.xml.temp | cut -d " " -f 1)" != "45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770" ]; then
    printf "root-anchors.xml SHA256 checksum mismatch !\n"
    echo "1" > health
    rm files/root-anchors.xml.temp
    return
  fi
  rm files/root-anchors.xml.updated
  mv files/root-anchors.xml.temp files/root-anchors.xml.updated
  KEYTAGS=$(xpath -q -e '/TrustAnchor/KeyDigest/KeyTag/node()' files/root-anchors.xml.updated)
  ALGORITHMS=$(xpath -q -e '/TrustAnchor/KeyDigest/Algorithm/node()' files/root-anchors.xml.updated)
  DIGESTTYPES=$(xpath -q -e '/TrustAnchor/KeyDigest/DigestType/node()' files/root-anchors.xml.updated)
  DIGESTS=$(xpath -q -e '/TrustAnchor/KeyDigest/Digest/node()' files/root-anchors.xml.updated)
  i=1
  rm files/root.key.updated
  while [ 1 ]; do
    KEYTAG=$(echo $KEYTAGS | cut -d" " -f$i)
    [ "$KEYTAG" != "" ] || break
    ALGORITHM=$(echo $ALGORITHMS | cut -d" " -f$i)
    DIGESTTYPE=$(echo $DIGESTTYPES | cut -d" " -f$i)
    DIGEST=$(echo $DIGESTS | cut -d" " -f$i)
    echo ". IN DS $KEYTAG $ALGORITHM $DIGESTTYPE $DIGEST" >> files/root.key.updated
    i=`expr $i + 1`
  done
  [ $VERBOSE = 0 ] || printf "Ran during $(($(date +%s)-$T_START)) seconds\n"
  unset -v T_START
  unset -v ALGORITHMS
  unset -v ALGORITHM
  unset -v DIGESTTYPES
  unset -v DIGESTTYPE
  unset -v DIGESTS
  unset -v DIGEST
}

main(){
  printf "\n =========================================\n"
  printf " =========================================\n"
  printf " ======= Updated Docker container ========\n"
  printf " =========================================\n"
  printf " =========================================\n"
  printf " == by github.com/qdm12 - Quentin McGaw ==\n\n"
  cd /updated || exit 1
  gitClone || exit 1
  while [ 1 ]; do
    gitPull || exit 1
    buildMaliciousHostnames || exit 1
    buildMaliciousIPs || exit 1
    #extendMaliciousIPs
    buildNamedRoot || exit 1
    buildRootAnchors || exit 1
    gitPush
    sleep 432000
  done
  status=$?
  printf "\n =========================================\n"
  printf " Exit with status $status\n"
  printf " =========================================\n"
}

main