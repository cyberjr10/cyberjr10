<?php

// Simulated Bitcoin transaction class
class BitcoinTransaction {
    public $recipientAddress;
    public $amount;
    public $expiryTime;

    public function __construct($recipientAddress, $amount, $expiryTime) {
        $this->recipientAddress = $recipientAddress;
        $this->amount = $amount;
        $this->expiryTime = $expiryTime;
    }
}

// Function to simulate sending Bitcoin
function sendBitcoin($recipientAddress, $amount, $expiryTime) {
    // Simulate sending Bitcoin
    $transaction = new BitcoinTransaction($recipientAddress, $amount, $expiryTime);
    return $transaction;
}

// Function to check if a transaction has expired
function isExpired($transaction) {
    // Check if the current time is greater than the expiry time
    return time() > $transaction->expiryTime;
}

// Check if the form has been submitted
if(isset($_POST['send_btc'])) {
    // Retrieve form data
    $recipientAddress = $_POST['recipient_address'];
    $amount = $_POST['btc_amount'];
    $expiryTime = time() + (48 * 60 * 60); // 48 hours expiry time

    // Send simulated Bitcoin transaction
    $transaction = sendBitcoin($recipientAddress, $amount, $expiryTime);

    // Check if the transaction has expired
    if (isExpired($transaction)) {
        echo "The transaction has expired.";
    } else {
        echo "Bitcoin sent to $recipientAddress: $amount BTC";
    }
}

?>
