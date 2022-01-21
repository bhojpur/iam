# Bhojpur IAM - Identity & Access Management

The Bhojpur IAM is used as an enterprise grade Identity &amp; Access Management system within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform). It features a centralized authentication and Single-Sign-On (SSO) platform based on OAuth 2.0 / OIDC standards.

## Quick Start

Run your own Bhojpur IAM program in a few minutes.

### Download

There are two methods, get source code via go subcommand `get`:

```shell
go get github.com/bhojpur/iam
```

  or `git`:

```bash
git clone https://github.com/bhojpur/iam.git
```

Finally, change the directory:

```bash
cd iam/
```

We provide two start up methods for all kinds of users.

### Manual

#### Simple configuration

The Bhojpur IAM requires a running relational database to be operational. Thus, you need to modify the configuration to point to your database instance.

Edit `conf/app.conf`, modify `dataSourceName` to correct database info, which follows this format:

```bash
username:password@tcp(database_ip:database_port)/
```

#### Run

The Bhojpur IAM provides two run modes, the difference is binary size and user prompt.

##### Development Mode

Edit the `conf/app.conf` file, set `runmode=dev`. Firstly, builds the front-end files:

```bash
cd pkg/webui && yarn && yarn run start

```
*The Bhojpur IAM's front-end module is built using yarn. You should use `yarn` instead of `npm`. It has a potential failure during building the files, if you use `npm`.*

Then, build the back-end binary files, change directory to root(Relative to the Bhojpur IAM):

```bash
go run server.go
go run client.go
```

That's it! Try to visit http://127.0.0.1:7001/. :small_airplane:  

**But make sure you always request the backend port 8000, when you are using SDKs.**

##### Production Mode

Edit the `conf/app.conf`, set `runmode=prod`. Firstly, build the front-end files:

```bash
cd webui/ && yarn && yarn run build
```

Then, build the back-end binary files, change directory to root(Relative to Bhojpur IAM):

```bash
go mod tidy
go get
go build -o bin/iamsvr server.go
go build -o bin/iamctl client.go
sudo bin/iamsvr
```

> Notice, you should visit back-end port, default 8000. Now try to visit **http://SERVER_IP:8000/**

### Docker

The Bhojpur IAM provides two kinds of binary images: 

- bhojpur/platform-iam-full, in which the Bhojpur IAM binary, a MySQL database, and all necessary configurations are packed up. This image is for new users to have a trial on the Bhojpur IAM quickly. **With this image you can start a Bhojpur IAM instance immediately with single command without any complex configuration**. **NOTE: we DO NOT recommend you to use this image in a productive environment**

- bhojpur/platform-iam-only: normal & graceful binary image with only Bhojpur IAM and environment installed. 

This method requires [docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/) to be installed first.

### Start Bhojpur IAM with bhojpur/platform-iam-full

if the image is not pulled, pull it from DockerHub

```shell
docker pull bhojpur/platform-iam-full
```

Start your instance with

```shell
docker run -p 8000:8000 bhojpur/platform-iam-full
```

Now, you can visit http://localhost:8000 and have a try. The default account's username and password is 'admin' and '123'. Go for it!

### Start the Bhojpur IAM with bhojpur/platform-iam-only

#### modify the configurations

For the convenience of your first attempt, docker-compose.yml contains commands to start a database via docker.

Thus, edit `conf/app.conf` to point to the location of your database(db:3306), modify `dataSourceName` to the fixed content:

```bash
dataSourceName = root:123456@tcp(db:3306)/
```

> If you need to modify `conf/app.conf`, you need to re-run `docker-compose up`.

#### Run

```bash
docker-compose up
```

That's it! Try to visit http://localhost:8000/. :small_airplane:

## Detailed documentation

We also provide complete product [documentation](https://docs.bhojpur.net/) as a reference.


## Contribute

For the Bhojpur IAM, if you have any questions, you can give Issues, or you can also directly start Pull Requests (but we recommend giving issues first to communicate with the community).

### I18N Translations

If you are contributing to the Bhojpur IAM, When you add some words in the ```pkg/webui/``` directory, please remember to add what you have added to the ```pkg/webui/src/locales/en/data.json``` file also.
