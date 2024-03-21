<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simulated Bitcoin Sender</title>
</head>
<body>
    <h1>Simulated Bitcoin Sender</h1>
    <form method="post">
        <label for="btc_amount">Bitcoin Amount:</label>
        <input type="number" id="btc_amount" name="btc_amount" step="0.0001" min="0.0001" required>
        <br>
        <label for="recipient_address">Recipient Address:</label>
        <input type="text" id="recipient_address" name="recipient_address" required>
        <br>
        <button type="submit" name="send_btc">Send Bitcoin</button>
    </form>
    <hr>

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
</body>
</html>
