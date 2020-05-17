# Word Stallion

This is a simple word game where you need to guess the right definition of a word to move your
horse to the finish line. The first horse to cross the finish line wins!

## Requirements

* AWS CLI already configured with Administrator permission
* [Golang](https://golang.org)

## Setup process

### Building

```shell
sam build
```

## Packaging and deployment

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

This assumes you own a top level domain on AWS. The deployment package will create a "wordstallion" subdomain for you.
