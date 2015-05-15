# RPC/TCP LogCabin protocol

## Version 1

Syntax:

    magic          ubint16 = 0xdaf4
    version        ubint16 = 1
    payload_length ubint32
    message_id     ubint64
    payload        [payload_length]byte

## Proposed version 2, encrypted transport

Still thinking; this probably is not a good design.

The syntax of the header is:

    payload_length ulint32
    timestamp      ulint32
    zeros          [8]byte
    
    eheader        [16]byte = AES-ECB(header)    
    nonce          [24]byte
    tag            [16]byte
    payload        [payload_length]byte

where payload is an AES-GCM encrypted message.

```py
magic = 'ad1b7beb'
digest = sha384('LogCabin RPC-over-TCP version 2').hexdigest()
print magic == digest[:8]
```
