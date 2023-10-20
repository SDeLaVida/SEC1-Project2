for ($i = 0; $i -le 3; $i++) {
    $x = $i
    Start-Process PowerShell "go run . $x | tee -filepath './logs/text$x.txt'"
}