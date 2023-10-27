#!/bin/bash
set -e

source $(dirname $0)/../base.sh

#${REGISTRY_IMAGE} dont work with https urls
gitRemoteUrl=$(git remote get-url origin)
if [[ $gitRemoteUrl == https* ]] ;then
    error "Pls use the ssh remote origin path: git@git.dev...."
    exit
fi

serviceName=$(jq -r .name scripts/service.json)
version=$(jq -r .version scripts/service.json)
namespaceName=$(jq -r .namespace scripts/service.json)
teamName=$(echo ${REGISTRY_IMAGE} | cut -d/ -f2)
servicePrefix=$(jq -r .prefix scripts/service.json)
registryImageRelativePath=$(echo ${REGISTRY_IMAGE} | cut -d/ -f2-)

read -p "Enter the ldap's of the codeOwners (e.g. @m.mustermann @t.tester) : " codeOwners

echo "The project will be created with the following values:"
echo "Service: ${serviceName}"
echo "Version: ${version}"
echo "Namespace: ${namespaceName}"
echo "Team: ${teamName}"
echo "CodeOwners: $codeOwners"
echo "Service Prefix: ${servicePrefix}"
echo "Registry Image: ${REGISTRY_IMAGE}"
echo "Relative Registry Path: ${registryImageRelativePath}"

read -p "Continue? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1

findCommand () {
    find . \
        ! -path '*/\.git*' \
        ! -path '*initProject.sh*' \
        ! -path '*initMonitoringCharts.sh*' \
        ! -path '*/\.idea*' \
        ! -path '*/\vendor*' \
        ! -path '*/\monitoring*' \
        ! -path '*README.md*' \
        ! -path '*.tar.gz' \
        ! -path '*.phar' \
    -type f -print0
}

findCommand | xargs -0 sed -i "s#{{your_service}}#${serviceName}#g"
findCommand | xargs -0 sed -i "s#{{your-namespace}}#${namespaceName}#g"
findCommand | xargs -0 sed -i "s/{{your_codeOwners}}/$codeOwners/g"
findCommand | xargs -0 sed -i "s#{{your-team}}#${teamName}#g"
findCommand | xargs -0 sed -i "s#{{your-prefix}}#${servicePrefix}#g"
findCommand | xargs -0 sed -i "s#{{your-registry-image}}#${REGISTRY_IMAGE}#g"
findCommand | xargs -0 sed -i "s#{{your-registry-relative-path}}#${registryImageRelativePath}#g"
findCommand | xargs -0 sed -i "s#{{your_serviceVersion}}#${version}#g"

