package a

type MyErr struct{}

func (MyErr) Error() string {
	return "My Error"
}

func f() *MyErr {
	return nil
}

func f2() error {
	var p *MyErr
	var err error

	for {
		if true {
			select {
			default:
				return p // want "Return pointer as error"
			}
		}
		break
	}
	func() *MyErr {
		return p
	}()

	return &MyErr{} // want "Return pointer as error"
	return p        // want "Return pointer as error"
	return f()      // want "Return pointer as error"
	return MyErr{}
	return err
	return nil
}

func f3() (error, error) {
	var err error
	var p *MyErr
	return nil, nil
	return err, err
	return nil, p // want "Return pointer as error"
	return p, nil // want "Return pointer as error"
}

func main() {
	var err error
	var p *MyErr

	err = &MyErr{} // want "Assign pointer to error"
	err = f()      // want "Assign pointer to error"
	err = p        // want "Assign pointer to error"
	err = MyErr{}
	err = nil
	err = (*MyErr)(nil) // want "Assign pointer to error"

	var i int
	var j int
	i, err, j = 0, p, 0 // want "Assign pointer to error"

	func() error {
		func() *MyErr {
			return nil
		}()
		return &MyErr{} // want "Return pointer as error"
	}()

	_, _, _ = i, j, err
}
