# Build guide

## Platform preparation

* GoLang 1.12+

    **Mac OSX**
    ```
    brew install go
    ```
    
* Python 3.7+ Virtual Environment

    **Mac OSX**
    ```
    brew install python
    pip install virtualenv setuptools wheel
    ```

## Environment

### Source checkout

First of all, you need to check out the source.
```bash
git clone $REPOSITORY_URL goloop
```

### Prepare virtual environment
```bash
cd $HOME/goloop
virtualenv -p python3 venv
. venv/bin/activate
```

### Install required packages
```bash
pip install -r pyee/requirements.txt
```

## Build

### Build executables

```bash
make
```

Output binaries are placed under `bin/` directory.


### Build python package

```bash
make pyexec
```

Output files are placed under `pyee/dist/` directory.

## Quick start

First step, you need to make a configuration for the node.

```bash
./bin/gochain --save_key_store wallet.json --save config.json
```

It stores generated configuration(`config.json`) along with wallet keystore
(`wallet.json`). If you don't specify any password, it uses `gochain` as 
password of the keystore. You may apply more options while it generates.

Now, you may start the server with it.

```bash
./bin/gochain -config config.json
```

You may send transaction with the wallet (`wallet.json`) for initial balance
of your wallet.

This is single node configuration. If you want to make a network with multiple
nodes, you need to make own genesis and node configurations.
