// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
// pada pakage consensus adalah sebuah implementasi yang berbeda dari consensus engine Ethereum.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainHeaderReader defines a small collection of methods needed to access the local
// blockchain during header verification.
// ChainHeaderReader interface adalah interface yang berisi beberapa method yang dibutuhkan untuk mengakses blockchain lokal saat verfikasi header.
type ChainHeaderReader interface {
	// Config retrieves the blockchain's chain configuration.
	// metoda 'config' akan mengembalikan konfigurasi rantai blockchain.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	// metoda 'current header' akan mengembalikan header saat ini dari rantai blockchain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	// metoda 'get header' akan mengembalikan header dari database berdasarkan hash dan nomor.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	// metoda 'get header by number' akan mengembalikan header dari database berdasarkan nomor.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	// metoda 'get header by hash' akan mengembalikan header dari database berdasarkan hash.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd retrieves the total difficulty from the database by hash and number.
	// metoda 'get td' akan mengembalikan total difficulty dari database berdasarkan hash dan nomor.
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
// 'ChainReader interface' adalah interface yang berisi koleksi methoda yang dibutuhkan untuk mengakses blockchain lokal saat verfikasi header dan/atau  verifikasi uncle.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock retrieves a block from the database by hash and number.
	// metoda 'get block' akan mengembalikan block dari database berdasarkan hash dan nomor.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine is an algorithm agnostic consensus engine.
// 'Engine interface' adalah antarmuka yang menyediakan algoritma agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	// metoda 'author' akan mengembalikan alamat Ethereum dari akun yang minted block tersebut, yang mungkin berbeda dari header coinbase jika consensus engine berbasis tanda tangan.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	// metoda 'verify header' akan mengecek apakah header sesuai dengan aturan consensus dari engine tertentu. 
	// Verifikasi tanda tangan dapat dilakukan secara opsional di sini, atau dapat dilakukan secara eksplisit melalui metoda 'verify seal'.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	// metoda 'verify headers' akan sama dengan metoda 'verify header', namun verifikasi header dalam batch secara bersamaan.
	// Metoda ini akan mengembalikan channel output untuk membatalkan operasi dan channel input untuk mengambil verifikasi secara asinkron (urutan adalah urutan input).
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	// metoda 'verify uncles' akan mengecek apakah block uncles sesuai dengan aturan consensus dari engine tertentu.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	// metoda 'prepare' akan menginisialisasi field consensus dari header block sesuai dengan aturan dari engine tertentu.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	// metoda 'finalize' akan menjalankan perubahan state post-transaksi (misalnya block rewards) tapi tidak mengbangunkan block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	// catatan : header block dan state database mungkin akan diperbarui untuk mengikuti aturan consensus yang terjadi di akhir (misalnya block rewards).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	// 
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	// catatan : header block dan state database mungkin akan diperbarui untuk mengikuti aturan consensus yang terjadi di akhir (misalnya block rewards).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	// metoda 'seal' akan menghasilkan permintaan baru untuk block input dan mengirimkan hasilnya ke channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	// catatan : metoda ini akan mengembalikan langsung dan akan mengirim hasil secara asinkron. Lebih dari satu hasil juga mungkin dikembalikan berdasarkan algoritma consensus.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	// metoda 'seal hash' akan mengembalikan hash dari block sebelumnya yang dibungkus.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	// metoda 'calc difficulty' adalah algoritma pengaturan level kesulitan. Metoda ini akan mengembalikan tingkat kesulitan yang baru.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	// metoda 'apis' akan mengembalikan RPC API yang disediakan oleh consensus engine ini.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close terminates any background threads maintained by the consensus engine.
	// metoda 'close' akan mengakhiri semua thread background yang ditangani oleh consensus engine.
	Close() error
}

// PoW is a consensus engine based on proof-of-work.
// PoW adalah consensus engine berbasis proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	// metoda 'hashrate' akan mengembalikan mining hashrate dari consensus engine PoW.
	Hashrate() float64
}
