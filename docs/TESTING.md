# Testing Guide

## Test Data
Sample test files are provided in `docs/testdata/` directory:

### Text Files
1. `sample.txt`
   ```
   Hello World
   This is a sample text file
   With multiple lines
   And some words to count
   ```
   Expected Results:
   - Lines: 5
   - Words: 12
   - Bytes: 65

### JSON Files
1. `sample.json`
   ```json
   {
     "name": "Test Data",
     "items": [
       {"id": 1, "value": "first"},
       {"id": 2, "value": "second"}
     ],
     "metadata": {
       "created": "2024-03-21",
       "version": "1.0"
     }
   }
   ```
   Expected Results:
   - Lines: 10
   - Words: 12
   - Bytes: 156

### CSV Files
1. `sample.csv`
   ```
   id,name,value
   1,first item,100
   2,second item,200
   3,third item,300
   ```
   Expected Results:
   - Lines: 4
   - Words: 12
   - Bytes: 89

### XML Files
1. `sample.xml`
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   <root>
     <items>
       <item id="1">First Item</item>
       <item id="2">Second Item</item>
     </items>
   </root>
   ```
   Expected Results:
   - Lines: 7
   - Words: 4
   - Bytes: 156

## Running Tests

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

### Test Coverage
```bash
go test -cover ./...
```

## Test Scenarios

### 1. File Analysis
```bash
./analyzer analyze docs/testdata/
```
Expected output:
```
INFO Processed docs/testdata/sample.txt: 5 lines, 12 words, 65 bytes in 0.1s
INFO Processed docs/testdata/sample.json: 10 lines, 12 words, 156 bytes in 0.1s
INFO Processed docs/testdata/sample.csv: 4 lines, 12 words, 89 bytes in 0.1s
INFO Processed docs/testdata/sample.xml: 7 lines, 4 words, 156 bytes in 0.1s
```

### 2. File Hashing
```bash
./analyzer hash docs/testdata/sample.txt
```
Expected output:
```
SHA256: <hash value>
```

### 3. Base64 Encoding
```bash
./analyzer encode docs/testdata/sample.txt
```
Expected output:
```
Base64: <encoded content>
```

### 4. Base64 Decoding
```bash
./analyzer decode <base64_content> output.txt
```
Expected output:
```
Decoded content written to: output.txt
```

## Error Cases to Test
1. Non-existent files
2. Invalid file formats
3. Empty files
4. Large files (>1MB)
5. Corrupted files
6. Permission issues
7. Invalid base64 content
8. Invalid output paths

## Performance Testing
1. Directory with 1000+ files
2. Files >100MB
3. Concurrent processing
4. Memory usage monitoring
5. Processing time benchmarks 