# Personal Finance Calculator

Calculator for personal fincances. Calculations are only accurate to tax residents of Victoria, Australia.

## Usage

```bash
# Calculate income tax
go run main.go tax

# Calculate first home options
go run main.go firsthome
```

You can create a `.env` file to skip having to enter your details every time you run the calculator.

```bash
cp .env.example .env
vim .env
```

Start the calculator with the following command:

## Features

#### Existing

- Add salary sacrifice to super.
- Calculate income tax.
- Calculate Medicare Levy.

#### Future

- Use tax rates (etc) from CSV files so they aren't hard-coded.
- Take region into account for tax rates, etc.
- Take private health cover into account for Medicare Levy Surcharge for salaries over $90k.
- Allow for one "profile" to calculate multiple income streams.
