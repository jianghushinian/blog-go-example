package examples

// NOTE: Go 1.23 之前迭代器现状
/*
// bufio.Reader
r := bufio.NewReader(...)
for {
    line, _ , err := r.ReadLine()
    if err != nil {
        break
    }
    // do something
}


// bufio.Scanner
scanner := bufio.NewScanner(...)
for scanner.Scan() {
    line := scanner.Text()
    // do something
}


// database/sql.Rows
rows, _ := db.QueryContext(...)
for rows.Next() {
    if err := rows.Scan(...); err != nil {
        break
    }
    // do something
}
*/
