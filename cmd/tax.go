package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tobyscott25/eofy-calc/helpers"
)

var CmdCalculateTax = &cobra.Command{
	Use:   "tax",
	Short: "Calculate tax",
	Long:  `A fast income tax calulator for Australians, written in Go by Toby Scott.`,
	Run: func(cmd *cobra.Command, args []string) {

		questions := []*survey.Question{
			{
				Name:   "GrossAnnualSalary",
				Prompt: &survey.Input{Message: "Your annual salary (eg. 65000 = $65,000)", Default: os.Getenv("PFC_GROSS_SALARY")},
			},
			{
				Name:   "SalarySacrificePercent",
				Prompt: &survey.Input{Message: "How much are you salary sacrificing? (eg. 10 = 10%)", Default: os.Getenv("PFC_SALARY_SACRIFICE_PERCENTAGE")},
			},
			{
				Name:   "HasHecsHelpDebt",
				Prompt: &survey.Confirm{Message: "Do you have a HECS/HELP debt?", Default: os.Getenv("PFC_HAS_HELP_DEBT") == "true"},
			},
			{
				Name:   "HasPrivHealthCover",
				Prompt: &survey.Confirm{Message: "Do you have private health insurance?", Default: os.Getenv("PFC_HAS_PRIVATE_HEALTH_COVER") == "true"},
			},
		}

		answers := struct {
			GrossAnnualSalary      int
			SalarySacrificePercent int
			HasHecsHelpDebt        bool
			HasPrivHealthCover     bool
		}{}

		err := survey.Ask(questions, &answers)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Hard-coded values, yet to take user input for these
		single := true // Does not have a spouse (married or de facto)
		// numberOfDependants := 0
		// saptoEligible := false // Were you eligible for the seniors and pensioners tax offset (SAPTO)?

		// Run calculations
		salarySacrificeAmount := answers.GrossAnnualSalary * answers.SalarySacrificePercent / 100
		taxableIncome := helpers.CalcTaxableIncome(answers.GrossAnnualSalary, answers.SalarySacrificePercent)
		incomeTax := helpers.CalcIncomeTax(float64(taxableIncome))
		hecsRepaymentRate := helpers.CalcHecsHelpRepaymentRate(answers.GrossAnnualSalary) // Repayment rate is based on your gross salary, not taxable income.
		hecsRepaymentAmount := float64(answers.GrossAnnualSalary) * hecsRepaymentRate
		paysMedicareLevySurcharge := helpers.PaysMedicareLevySurcharge(float64(taxableIncome), single, answers.HasPrivHealthCover)
		medicareLevy := helpers.CalcMedicareLevy(float64(taxableIncome), single, paysMedicareLevySurcharge)

		// Print out report
		fmt.Println("=====================================")
		fmt.Printf("Your gross annual salary is $%d\n", answers.GrossAnnualSalary)
		fmt.Printf("You are salary sacrificing $%d (%d%%)\n", salarySacrificeAmount, answers.SalarySacrificePercent)
		fmt.Printf("Your taxable income is $%d\n", taxableIncome)
		fmt.Println("=====================================")
		fmt.Printf("Your income tax is $%.2f\n", incomeTax)
		fmt.Printf("Your medicare levy is $%.2f\n", medicareLevy)
		if answers.HasHecsHelpDebt {
			fmt.Printf("You have a HELP/HECS debt, your HECS repayment is $%.2f (%.2f%% of $%d)\n", hecsRepaymentAmount, (hecsRepaymentRate * 100), answers.GrossAnnualSalary)
		}
		fmt.Println("=====================================")
		fmt.Printf("Total going to the ATO: $%.2f\n", (incomeTax + medicareLevy + hecsRepaymentAmount))
		fmt.Printf("Salary sacrifice amount: $%d\n", salarySacrificeAmount)
		fmt.Printf("Total take home pay: $%.2f\n", (float64(answers.GrossAnnualSalary) - (incomeTax + medicareLevy + hecsRepaymentAmount) - float64(salarySacrificeAmount)))
		fmt.Println("=====================================")
	},
}
