package api

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"github.com/w33ked/theblockchain/core"
	"github.com/w33ked/theblockchain/types"
)

type APIError struct {
	Error string
}

type Block struct {
	Hash          string
	Version       uint32
	DataHash      string
	PrevBlockHash string
	Height        uint32
	Timestamp     int64
	Validator     string
	Signature     string
}

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	ServerConfig
	bc *core.Blockchain
}

func NewServer(cfg ServerConfig, bc *core.Blockchain) *Server {
	return &Server{
		ServerConfig: cfg,
		bc:           bc,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrId := c.Param("hashorid")
	height, err := strconv.Atoi(hashOrId)
	// if err is nil we can assume the height of the block is given
	if err == nil {
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, intoJSONBlock(block))
	}

	// oherwise assume it is the hash
	b, err := hex.DecodeString(hashOrId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	block, err := s.bc.GetBlockByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func intoJSONBlock(block *core.Block) Block {
	return Block{
		Hash:          block.Hash(core.BlockHasher{}).String(),
		Version:       block.Version,
		Height:        block.Header.Height,
		DataHash:      block.Header.DataHash.String(),
		PrevBlockHash: block.DataHash.String(),
		Timestamp:     block.Header.Timestamp,
		Validator:     block.Validator.Address().String(),
		Signature:     block.Signature.String(),
	}
}
