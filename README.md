# gochain
A simple blockchain implementation in Go. Inspired and guided by [Jeiwan's Blog](https://jeiwan.net/)


### Block
In blockchain it’s blocks that store valuable information. For example, bitcoin blocks store transactions, the essence of any cryptocurrency. Besides this, a block contains some technical information, like its version, current timestamp and the hash of the previous block.

### Hashing
The way hashes are calculates is very important feature of blockchain, and it’s this feature that makes blockchain secure. The thing is that calculating a hash is a computationally difficult operation, it takes some time even on fast computers (that’s why people buy powerful GPUs to mine Bitcoin). This is an intentional architectural design, which makes adding new blocks difficult, thus preventing their modification after they’re added

Bitcoin mining uses the hashcash proof of work function; the hashcash algorithm requires the following parameters: a service string, a nonce, and a counter. In bitcoin the service string is encoded in the block header data structure, and includes a version field, the hash of the previous block, the root hash of the merkle tree of all transactions in the block, the current time, and the difficulty. Bitcoin stores the nonce in the extraNonce field which is part of the coinbase transaction, which is stored as the left most leaf node in the merkle tree (the coinbase is the special first transaction in the block). The counter parameter is small at 32-bits so each time it wraps the extraNonce field must be incremented (or otherwise changed) to avoid repeating work.
[Details here](https://en.bitcoin.it/wiki/Block_hashing_algorithm)

### Blockchain
In its essence blockchain is just a database with certain structure: it’s an ordered, back-linked list. Which means that blocks are stored in the insertion order and that each block is linked to the previous one. This structure allows to quickly get the latest block in a chain and to (efficiently) get a block by its hash.
 
### Proof of Work
This mechanism is very similar to the one from real life: one has to work hard to get a reward and to sustain their life. In blockchain, some participants (miners) of the network work to sustain the network, to add new blocks to it, and get a reward for their work. As a result of their work, a block is incorporated into the blockchain in a secure way, which maintains the stability of the whole blockchain database. It’s worth noting that, the one who finished the work has to prove this.

Proof-of-Work algorithms must meet a requirement: doing the work is hard, but verifying the proof is easy. A proof is usually handed to someone else, so for them, it shouldn’t take much time to verify it.

### Transactions
In Bitcoin, payments are realized in completely different way. There are:
 * No accounts.
 * No balances.
 * No addresses.
 * No coins.
 * No senders and receivers.

A transaction is a transfer of Bitcoin value that is broadcast to the network and collected into blocks. A transaction typically references previous transaction outputs as new transaction inputs and dedicates all input Bitcoin values to new outputs. Transactions are not encrypted, so it is possible to browse and view every transaction ever collected into a block. Once transactions are buried under enough confirmations they can be considered irreversible. 
[Read More here](https://en.bitcoin.it/wiki/Transaction)
