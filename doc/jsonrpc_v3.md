---
title: JSON-RPC v3
---

# Goloop JSON-RPC API v3

## Introduction

This document explains JSON-RPC APIs (version 3) available to interact with Goloop nodes.

## Value Types

Basically, every VALUE in JSON-RPC message is string.
Below table shows the most common "VALUE types".

| VALUE type | Description | Example |
|:----------|:----|:----|
| <a id="T_ADDR_EOA">T_ADDR_EOA</a> | "hx" + 40 digit HEX string | hxbe258ceb872e08851f1f59694dac2558708ece11 |
| <a id="T_ADDR_SCORE">T_ADDR_SCORE</a> | "cx" + 40 digit HEX string | cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32 |
| <a id="T_HASH">T_HASH</a> | "0x" + 64 digit HEX string | 0xc71303ef8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238 |
| <a id="T_INT">T_INT</a> | "0x" + lowercase HEX string | 0xa |
| <a id="T_BIN_DATA">T_BIN_DATA</a> | "0x" + lowercase HEX string. Length must be even. | 0x34b2 |
| <a id="T_SIG">T_SIG</a> | base64 encoded string | VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA= |
| <a id="T_DATA_TYPE">T_DATA_TYPE</a> | Type of data | call, deploy, or message |

## Error Codes

This chapter explains the error codes used in Goloop JSON-RPC API response.

Below table shows the default error messages for the error code. Actual message may vary depending on the implementation.

| Category   | Error code | Message | Description |
|:-----------|:---------|:------|:-----|
|Json Parsing| -32700 | Parse error | Invalid JSON was received by the server.<br/>An error occurred on the server while parsing the JSON text. |
|RPC Parsing | -32600 | Invalid Request | The JSON sent is not a valid Request object. |
|            | -32601 | Method not found | The method does not exist / is not available. |
|            | -32602 | Invalid params | Invalid method parameter(s). |
|            | -32603 | Internal error | Internal JSON-RPC error. |
|Server Error| -32000 ~ -32099 |  | Server error. |
|System Error| -31000 | System Error | Unknown system error. |
|            | -31001 | Pool Overflow | Transaction pool overflow. |
|            | -31002 | Pending | Transaction is in the pool, but not included in the block. |
|            | -31003 | Executing | Transaction is included in the block, but it doesn’t have confirmed result. |
|            | -31004 | Not found | Requested data is not found. |
|SCORE Error | -30000 ~ -30099 |  | Mapped error codes from Core2 Design - SCORE Result. |

## JSON-RPC Methods

### icx_getLastBlock

Returns the last block information.

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_getLastBlock",
}
```
#### Parameters

None

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": {
    "block_hash": "8e25acc5b5c74375079d51828760821fc6f54283656620b1d5a715edcc0770c6",
     "confirmed_transaction_list": [
      {
        "from": "hx84f6c686fba03bc7ca65d15ae844ee56ff24a32b",
        "nid": "0x1",
        "signature": "tCUwOb6vsaUKy+NYvmzdJYC0jm3Erd5cR6wKnVuAjzMOECC+t/oK7fG/Tz2Y3C25o0AfCmbneXpias6xco+43wE=",
        "stepLimit": "0x3e8",
        "timestamp": "0x58a14bfe9b904",
        "to": "hx244deea00413d85c6637e7fdd53afa697f29d08f",
        "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f",
        "value": "0xa",
        "version": "0x3"
      }
     ],
    "height": 512,
    "merkle_tree_root_hash": "5c8d4e59ded657c6acbb67030929dfcaf114a268d6d58df53e7174e40db74158",
    "peer_id": "hx4208599c8f58fed475db747504a80a311a3af63b",
    "prev_block_hash": "0fdf04d13229482e3533948d4582344a3d44c399e71ab12c653ae57bcbee5d90",
    "signature": "",
    "time_stamp": 1559204699330360,
    "version": "2.0"
  },
  "id": "1001"
}
```
#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

### icx_getBlockByHeight

Returns block information by block height.

> Request

```json
{
  "id": "1001",
  "jsonrpc": "2.0",
  "method": "icx_getBlockByHeight",
  "params": {
    "height": "0x100"
  }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:-----------|:-----|
| height | [T_INT](#T_INT) | Integer of a block height |

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": {
    "block_hash": "8e25acc5b5c74375079d51828760821fc6f54283656620b1d5a715edcc0770c6",
     "confirmed_transaction_list": [
      {
        "from": "hx84f6c686fba03bc7ca65d15ae844ee56ff24a32b",
        "nid": "0x1",
        "signature": "tCUwOb6vsaUKy+NYvmzdJYC0jm3Erd5cR6wKnVuAjzMOECC+t/oK7fG/Tz2Y3C25o0AfCmbneXpias6xco+43wE=",
        "stepLimit": "0x3e8",
        "timestamp": "0x58a14bfe9b904",
        "to": "hx244deea00413d85c6637e7fdd53afa697f29d08f",
        "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f",
        "value": "0xa",
        "version": "0x3"
      }
     ],
    "height": 512,
    "merkle_tree_root_hash": "5c8d4e59ded657c6acbb67030929dfcaf114a268d6d58df53e7174e40db74158",
    "peer_id": "hx4208599c8f58fed475db747504a80a311a3af63b",
    "prev_block_hash": "0fdf04d13229482e3533948d4582344a3d44c399e71ab12c653ae57bcbee5d90",
    "signature": "",
    "time_stamp": 1559204699330360,
    "version": "2.0"
  },
  "id": "1001"
}
```

#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

### icx_getBlockByHash

Returns block information by block hash.

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_getBlockByHeight",
  "params": {
      "hash": "8e25acc5b5c74375079d51828760821fc6f54283656620b1d5a715edcc0770c6"
  }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:-----------|:-----|
| hash | [T_HASH](#T_HASH) | Hash of a block |

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": {
    "block_hash": "8e25acc5b5c74375079d51828760821fc6f54283656620b1d5a715edcc0770c6",
     "confirmed_transaction_list": [
      {
        "from": "hx84f6c686fba03bc7ca65d15ae844ee56ff24a32b",
        "nid": "0x1",
        "signature": "tCUwOb6vsaUKy+NYvmzdJYC0jm3Erd5cR6wKnVuAjzMOECC+t/oK7fG/Tz2Y3C25o0AfCmbneXpias6xco+43wE=",
        "stepLimit": "0x3e8",
        "timestamp": "0x58a14bfe9b904",
        "to": "hx244deea00413d85c6637e7fdd53afa697f29d08f",
        "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f",
        "value": "0xa",
        "version": "0x3"
      }
     ],
    "height": 512,
    "merkle_tree_root_hash": "5c8d4e59ded657c6acbb67030929dfcaf114a268d6d58df53e7174e40db74158",
    "peer_id": "hx4208599c8f58fed475db747504a80a311a3af63b",
    "prev_block_hash": "0fdf04d13229482e3533948d4582344a3d44c399e71ab12c653ae57bcbee5d90",
    "signature": "",
    "time_stamp": 1559204699330360,
    "version": "2.0"
  },
  "id": "1001"
}
```

#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

### icx_call

Calls SCORE's external function.

Does not make state transition (i.e., read-only).

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_call",
  "params": {
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11", // TX sender address
        "to": "cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32",   // SCORE address
        "dataType": "call",
        "data": {
            "method": "get_balance", // SCORE external function
            "params": {
                "address": "hx1f9a3310f60a03934b917509c86442db703cbd52" // input parameter of "get_balance"
            }
        }
    }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:-----------|:------------|
| from | [T_ADDR_EOA](#T_ADDR_EOA) | Message sender's address. |
| to | [T_ADDR_SCORE](#T_ADDR_SCORE) | SCORE address that will handle the message. |
| dataType | [T_DATA_TYPE](#T_DATA_TYPE) | `call` is the only possible data type. |
| data | T_DICT | See [Parameters - data](#sendtxparameterdata). |
| data.method | String | Name of the function. |
| data.params | T_DICT | Parameters to be passed to the function. |

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": "0x2961fff8ca4a62327800000",
  "id": 1001
}
```

#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success||

### icx_getBalance

Returns the ICX balance of the given EOA or SCORE.

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_getBalance",
   "params": {
        "address": "hxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32"
    }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:-----------|:-----|
| address | [T_ADDR_EOA](#T_ADDR_EOA) or [T_ADDR_SCORE](#T_ADDR_SCORE) | Address of EOA or SCORE |

> Example responses

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "result": "0xde0b6b3a7640000"
}
```
#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success||

### icx_getScoreApi

Returns SCORE's external API list.

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_getScoreApi",
  "params": {
      "address": "cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32"  // SCORE address
  }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:-----------|:-----|
| address | [T_ADDR_SCORE](#T_ADDR_SCORE) | SCORE adress to be examined. |

> Example responses

```json
{
    "jsonrpc": "2.0",
    "id": 1234,
    "result": [
        {
            "type": "function",
            "name": "balanceOf",
            "inputs": [
                {
                    "name": "_owner",
                    "type": "Address"
                }
            ],
            "outputs": [
                {
                    "type": "int"
                }
            ],
            "readonly": "0x1"
        },
        {
            "type": "eventlog",
            "name": "FundTransfer",
            "inputs": [
                {
                    "name": "backer",
                    "type": "Address",
                    "indexed": "0x1"
                },
                {
                    "name": "amount",
                    "type": "int",
                    "indexed": "0x1"
                },
                {
                    "name": "is_contribution",
                    "type": "bool",
                    "indexed": "0x1"
                }
            ]
        },
        {...}
    ]
}
```
#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success||

* Fields containing information about the function
    - type : `function`, `fallback`, or `eventlog`
    - name : function name
    - inputs : parameters in array
        + name : parameter name
        + type : parameter type (`int`, `str`, `bytes`, `bool`, `Address`)
        + indexed : `0x1` if the parameter is indexed (when this is `eventlog`)
    - outputs : return value
        + type : return value type (`int`, `str`, `bytes`, `bool`, `Address`, `dict`, `list`)
    - readonly : `0x1` if this is declared as `external(readonly=True)`
    - payable : `0x1` if this has `payable` decorator

### icx_getTotalSupply

Returns total ICX coin supply that has been issued.

> Request

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "method": "icx_getTotalSupply",
}
```
#### Parameters

None

> Example responses

```json
{
  "id": 1001,
  "jsonrpc": "2.0",
  "result": "0x2961fff8ca4a62327800000"
}
```

#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success||

### icx_getTransactionResult

Returns the transaction result requested by transaction hash.

> Request

```json
{
  "jsonrpc": "2.0",
  "id": "1001",
  "method": "icx_getTransactionResult",
  "params": {
    "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f"
  }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:----------|:----- |
| txHash | [T_HASH](#T_HASH) | Hash of the transaction |

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": {
    "blockHash": "0x8ef3b2a67262b9b1fe4b598059774472e9ccef401734335d87a4ba998cfd40fb",
    "blockHeight": "0x200",
    "cumulativeStepUsed": "0x0",
    "eventLogs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1",
    "stepPrice": "0x0",
    "stepUsed": "0x0",
    "to": "hx244deea00413d85c6637e7fdd53afa697f29d08f",
    "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f",
    "txIndex": "0x0"
  },
  "id": "1001"
}
```
#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

| KEY | VALUE type | Description |
|:----|:----------|:-----|
| status | [T_INT](#T_INT) | 1 on success, 0 on failure. |
| to | [T_ADDR_EOA](#T_ADDR_EOA) or [T_ADDR_SCORE](#T_ADDR_SCORE) | Recipient address of the transaction |
| failure | T_DICT | This field exists when status is 0. Contains code(str) and message(str). |
| txHash | [T_HASH](#T_HASH) | Transaction hash |
| txIndex | [T_INT](#T_INT) | Transaction index in the block |
| blockHeight | [T_INT](#T_INT) | Height of the block that includes the transaction. |
| blockHash | [T_HASH](#T_HASH) | Hash of the block that includes the transation. |
| cumulativeStepUsed | [T_INT](#T_INT) | Sum of stepUsed by this transaction and all preceeding transactions in the same block. |
| stepUsed | [T_INT](#T_INT) | The amount of step used by this transaction. |
| stepPrice | [T_INT](#T_INT) | The step price used by this transaction. |
| scoreAddress | [T_ADDR_SCORE](#T_ADDR_SCORE) | SCORE address if the transaction created a new SCORE. (optional) |
| eventLogs | [T_ARRAY](#T_ARRAY) | Array of eventlogs, which this transaction generated. |
| logsBloom | [T_BIN_DATA](#T_BIN_DATA) | Bloom filter to quickly retrieve related eventlogs. |

### icx_getTransactionByHash

Returns the transaction information requested by transaction hash.

> Request

```json
{
  "jsonrpc": "2.0",
  "id": "1001",
  "method": "icx_getTransactionByHash",
  "params": {
    "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f"
  }
}
```
#### Parameters

| KEY | VALUE type | Description |
|:----|:----------|:----- |
| txHash | [T_HASH](#T_HASH) | Hash of the transaction |

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": {
    "blockHash": "0x8ef3b2a67262b9b1fe4b598059774472e9ccef401734335d87a4ba998cfd40fb",
    "blockHeight": "0x200",
    "from": "hx84f6c686fba03bc7ca65d15ae844ee56ff24a32b",
    "nid": "0x1",
    "signature": "tCUwOb6vsaUKy+NYvmzdJYC0jm3Erd5cR6wKnVuAjzMOECC+t/oK7fG/Tz2Y3C25o0AfCmbneXpias6xco+43wE=",
    "stepLimit": "0x3e8",
    "timestamp": "0x58a14bfe9b904",
    "to": "hx244deea00413d85c6637e7fdd53afa697f29d08f",
    "txHash": "0xd8da71e926052b960def61c64f325412772f8e986f888685bc87c0bc046c2d9f",
    "txIndex": "0x0",
    "value": "0xa",
    "version": "0x3"
  },
  "id": "1001"
}
```
#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

| KEY | VALUE type | Description |
|:----|:----------|:-----|
| version | [T_INT](#T_INT) | Protocol version ("0x3" for V3) |
| from | [T_ADDR_EOA](#T_ADDR_EOA) | EOA address that created the transaction |
| to | [T_ADDR_EOA](#T_ADDR_EOA) or [T_ADDR_SCORE](#T_ADDR_SCORE) | EOA address to receive coins, or SCORE address to execute the transaction. |
| value | [T_INT](#T_INT) | Amount of ICX coins in loop to transfer. When ommitted, assumes 0. (1 icx = 1 ^ 18 loop) |
| stepLimit |[T_INT](#T_INT) | Maximum step allowance that can be used by the transaction. |
| timestamp | [T_INT](#T_INT) | Transaction creation time. timestamp is in microsecond. |
| nid | [T_INT](#T_INT) |  Network ID |
| nonce | [T_INT](#T_INT) | An arbitrary number used to prevent transaction hash collision. |
| txHash | [T_HASH](#T_HASH) | Transaction hash |
| txIndex | [T_INT](#T_INT) | Transaction index in a block. Null when it is pending. |
| blockHeight | [T_INT](#T_INT) | Block height where this transaction was in. Null when it is pending. |
| blockHash | [T_HASH](#T_HASH) | Hash of the block where this transaction was in. Null when it is pending. |
| signature | [T_SIG](#T_SIG) | Signature of the transaction. |
| dataType | [T_DATA_TYPE](#T_DATA_TYPE) | Type of data. (call, deploy, or message) |
| data | T_DICT or String | Contains various type of data depending on the dataType. See [Parameters - data](#sendtxparameterdata). |

### icx_sendTransaction

You can do one of the followings using this function.
* Transfer designated amount of ICX coins from 'from' address to 'to' address.
* Install a new SCORE.
* Update the SCORE in the 'to' address.
* Invoke a function of the SCORE in the 'to' address.
* Transfer a message.

This function causes state transition.

> Coin transfer

```json
{
    "jsonrpc": "2.0",
    "method": "icx_sendTransaction",
    "id": 1234,
    "params": {
        "version": "0x3",
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "to": "hx5bfdb090f43a808005ffc27c25b213145e80b7cd",
        "value": "0xde0b6b3a7640000",
        "stepLimit": "0x12345",
        "timestamp": "0x563a6cf330136",
        "nid": "0x3",
        "nonce": "0x1",
        "signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA="
    }
}
```

> SCORE function call

```json
{
    "jsonrpc": "2.0",
    "method": "icx_sendTransaction",
    "id": 1234,
    "params": {
        "version": "0x3",
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "to": "cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32",
        "stepLimit": "0x12345",
        "timestamp": "0x563a6cf330136",
        "nid": "0x3",
        "nonce": "0x1",
        "signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA=",
        "dataType": "call",
        "data": {
            "method": "transfer",
            "params": {
                "to": "hxab2d8215eab14bc6bdd8bfb2c8151257032ecd8b",
                "value": "0x1"
            }
        }
    }
}
```

> SCORE install

```json
{
    "jsonrpc": "2.0",
    "method": "icx_sendTransaction",
    "id": 1234,
    "params": {
        "version": "0x3",
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "to": "cx0000000000000000000000000000000000000000", // address 0 means SCORE install
        "stepLimit": "0x12345",
        "timestamp": "0x563a6cf330136",
        "nid": "0x3",
        "nonce": "0x1",
        "signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA=",
        "dataType": "deploy",
        "data": {
            "contentType": "application/zip",
            "content": "0x1867291283973610982301923812873419826abcdef91827319263187263a7326e...", // compressed SCORE data
            "params": {  // parameters to be passed to on_install()
                "name": "ABCToken",
                "symbol": "abc",
                "decimals": "0x12"
            }
        }
    }
}
```

> SCORE update

```json
{
    "jsonrpc": "2.0",
    "method": "icx_sendTransaction",
    "id": 1234,
    "params": {
        "version": "0x3",
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "to": "cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32", // SCORE address to be updated
        "stepLimit": "0x12345",
        "timestamp": "0x563a6cf330136",
        "nid": "0x3",
        "nonce": "0x1",
        "signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA=",
        "dataType": "deploy",
        "data": {
            "contentType": "application/zip",
            "content": "0x1867291283973610982301923812873419826abcdef91827319263187263a7326e...", // compressed SCORE data
            "params": {  // parameters to be passed to on_update()
                "amount": "0x1234"
            }
        }
    }
}
```

> Message transfer

```json
{
    "jsonrpc": "2.0",
    "method": "icx_sendTransaction",
    "id": 1234,
    "params": {
        "version": "0x3",
        "from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "to": "hxbe258ceb872e08851f1f59694dac2558708ece11",
        "stepLimit": "0x12345",
        "timestamp": "0x563a6cf330136",
        "nid": "0x3",
        "nonce": "0x1",
        "signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA=",
        "dataType": "message",
        "data": "0x4c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c69742c2073656420646f20656975736d6f642074656d706f7220696e6369646964756e74207574206c61626f726520657420646f6c6f7265206d61676e6120616c697175612e20557420656e696d206164206d696e696d2076656e69616d2c2071756973206e6f737472756420657865726369746174696f6e20756c6c616d636f206c61626f726973206e69736920757420616c697175697020657820656120636f6d6d6f646f20636f6e7365717561742e2044756973206175746520697275726520646f6c6f7220696e20726570726568656e646572697420696e20766f6c7570746174652076656c697420657373652063696c6c756d20646f6c6f726520657520667567696174206e756c6c612070617269617475722e204578636570746575722073696e74206f6363616563617420637570696461746174206e6f6e2070726f6964656e742c2073756e7420696e2063756c706120717569206f666669636961206465736572756e74206d6f6c6c697420616e696d20696420657374206c61626f72756d2e"
    }
}
```

#### Parameters

| KEY | VALUE type | Required | Description |
|:----|:----------|:----:|:-----|
| version | [T_INT](#T_INT) | required | Protocol version ("0x3" for V3) |
| from | [T_ADDR_EOA](#T_ADDR_EOA) | required | EOA address that created the transaction |
| to | [T_ADDR_EOA](#T_ADDR_EOA) or [T_ADDR_SCORE](#T_ADDR_SCORE) | required | EOA address to receive coins, or SCORE address to execute the transaction. |
| value | [T_INT](#T_INT) | optional | Amount of ICX coins in loop to transfer. When ommitted, assumes 0. (1 icx = 1 ^ 18 loop) |
| stepLimit |[T_INT](#T_INT) | required | Maximum step allowance that can be used by the transaction. |
| timestamp | [T_INT](#T_INT) | required | Transaction creation time. timestamp is in microsecond. |
| nid | [T_INT](#T_INT) | required | Network ID ("0x1" for Mainnet, "0x2" for Testnet, etc) |
| nonce | [T_INT](#T_INT) | optional | An arbitrary number used to prevent transaction hash collision. |
| signature | [T_SIG](#T_SIG) | required | Signature of the transaction. |
| dataType | [T_DATA_TYPE](#T_DATA_TYPE) | optional | Type of data. (call, deploy, or message) |
| data | T_DICT or String | optional | The content of data varies depending on the dataType. See [Parameters - data](#sendtxparameterdata). |

#### <a id ="sendtxparameterdata">Parameters - data</a>
`data` contains the following data in various formats depending on the dataType.

#### dataType == call

It is used when calling a function in SCORE, and `data` has dictionary value as follows.

| KEY | VALUE type | Required | Description |
|:----|:-----------|:--------:|:------------|
| method | String | required | Name of the function to invoke in SCORE |
| params | T_DICT | optional | Function parameters |

##### dataType == deploy

It is used when installing or updating a SCORE, and `data` has dictionary value as follows.

| KEY | VALUE type | Required | Description |
|:----|:-----------|:--------:|:------------|
| contentType | String | required | Mime-type of the content |
| content | [T_BIN_DATA](#T_BIN_DATA) | required | Compressed SCORE data |
| params | T_DICT | optional | Function parameters will be delivered to on_install() or on_update() |

##### dataType == message

It is used when transfering a message, and `data` has a HEX string.

> Example responses

```json
{
  "jsonrpc": "2.0",
  "result": "0x402b630c5ed80d1b8f0d89ca14a091084bcc0f6a98bc52329bccc045415bc0bd",
  "id": "1001"
}
```

#### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|OK|Success|Block|

* Transaction hash ([T_HASH](#T_HASH)) on success
* Error code and message on failure