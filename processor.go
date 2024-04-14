// CommandProcessor interface implementation for RESP.
// Invoked by the network layer
package main

type RespCommandProcessor struct{}

func newRespCommandProcessor() *RespCommandProcessor {
	return &RespCommandProcessor{}
}

func (rcp RespCommandProcessor) processCommand(bytes []byte) []byte {

	_, err := newRespRequest(bytes)

	if err != nil {
		response := newRespResponse(RESP_SIMPLE_ERROR, []string{RESP_ERR})
		return response.marshalToBytes()
	}

	return nil
}
