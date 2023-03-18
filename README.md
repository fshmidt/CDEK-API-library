# CDEK API library

<!-- ToC start -->
# Content

1. [Description](#Description)
2. [Realisation](#Realisation)
3. [Installation](#Installation)
4. [Usage](#Usage)

[//]: # (1. [Примеры]&#40;#Примеры&#41;)
<!-- ToC end -->

# Description

This is a Golang library for working with the CDEK API. It allows you to easily calculate shipping prices based on package size and delivery addresses.

# Realisation

- Bearer Auth
- Clean Architecture

# Installation

To install the library, use the following command:

```go get github.com/fshmit/CDEK-API-library```

# Usage

To use the library, first import it into your Go project:

```import "github.com/fshmit/CDEK-API-library"```

Next, create a new client with your CDEK account credentials:

```
client, err := cdek_api_lib.NewClient(username, password, testMode, apiAddress)
if err != nil {
// Handle error
}
```

The testMode parameter specifies whether to use CDEK's test API or the live API. The apiAddress parameter is the URL of the CDEK API.

You can then use the client to calculate shipping prices:

```
prices, err := client.Calculate(fromAddress, toAddress, packageSize)
if err != nil {
// Handle error
}

for _, price := range prices {
// Process price information
}
```

The fromAddress and toAddress parameters should be the delivery addresses in string format. The packageSize parameter should be a CDEK_API_lib.Size struct that specifies the dimensions and weight of the package.
