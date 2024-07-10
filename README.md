[![Test](https://github.com/robbert229/pggen/workflows/Test/badge.svg)](https://github.com/robbert229/pggen/actions?query=workflow%3ATest) 
[![Lint](https://github.com/robbert229/pggen/workflows/Lint/badge.svg)](https://github.com/robbert229/pggen/actions?query=workflow%3ALint) 
[![GoReportCard](https://goreportcard.com/badge/github.com/robbert229/pggen)](https://goreportcard.com/report/github.com/robbert229/pggen)

# Notice

This is a temporary fork that I have been working on with the intention of 
backporting the changes to the upstream pggen. It is not fully functional.

# Changes

Because of the changes in how types are managed in v4 compared to v5, type registration needs to be configured in the AfterConnect function on the pgx config.

```golang
dbconfig, err := pgxpool.ParseConfig(databaseURL)
if err != nil {
	// handle error
}

dbconfig.AfterConnect = func(ctx context.Context, pgconn *pgconn.PgConn) error {
    err := your_gen_pkg.Register(ctx, pgconn)
    if err != nil {
        return fmt.Errorf("failed to register types: %w", err)
    }

    return nil
}
```

# Features Supported

The following examples compile, and are implemented.

* example/author
* example/domain
* example/go_pointer_types
* example/inline_param_count
* example/pgcrypto
* example/slices
* example/syntax
* example/custom_types
* example/enums
* example/device
* example/function
* example/separate_out_dir/out
* example/erp/order
* example/ltree
* example/void
* example/composite

# Feature Unsupported
* example/citext
* example/nested
* example/complex_params
