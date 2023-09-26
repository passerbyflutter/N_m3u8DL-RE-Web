#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

output_name=$2
if [[ -z "$output_name" ]]; then
  echo "usage: $0 <package-name> <output_name>"
  exit 1
fi
package_split=(${package//\// })
    
platforms=("windows/amd64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_file_name=$output_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi    

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$RELEASE_VERSION -X main.buildTime=$NOW -X main.isRelease=true" -o $output_file_name $package
    if [ $? -ne 0 ]; then
           echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done