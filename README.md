
[![Build Status](https://travis-ci.org/jfrey/terraform-provider-rackhd.svg?branch=master)](https://travis-ci.org/jfrey/terraform-provider-rackhd)

Prior to running these steps install Terraform (brew install terraform).

1. go build -o terraform-provider-rackhd
2. cp ./terraform-provider-rackhd $(dirname `which terraform`)
3. terraform plan
