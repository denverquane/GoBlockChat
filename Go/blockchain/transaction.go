package blockchain

type TransType string

const (
	ADD_MESSAGE TransType = "ADD_MESSAGE"
	DELETE_MESSAGE = "DELETE_MESSAGE"
	ADD_USER = "ADD_USER"
)

type AuthTransaction struct {
	Username 		string
	Password		string
	Channel 		string
	Message 		string
	TransactionType	string
}

type Transaction struct {
	Username        string
	Channel         string
	Message         string
	TransactionType string
}

func (trans AuthTransaction) RemovePassword() Transaction {
	return Transaction{Username:trans.Username, Channel:trans.Channel,
		Message:trans.Message, TransactionType:trans.TransactionType}
}

func (trans Transaction) ToString() string {
	return trans.Username + " posted \"" + trans.Message + "\" on the " + trans.Channel + " channel"
}

func SampleTransaction() Transaction {
	return Transaction{"samplehashhere", "Test", "Sample message.", "ADD_MESSAGE"}
}

func GetTransactionFormat() string {
	return "{username:User,password:Pass,channel:TestChannel,message:SampleMessage,transactiontype:ADD_MESSAGE}"
}