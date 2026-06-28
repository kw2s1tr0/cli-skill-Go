package controller

import (
	"aiagentcliapp/repository"
	repositorydepartment "aiagentcliapp/repository/department"
	repositoryemployee "aiagentcliapp/repository/employee"
	repositorylogin "aiagentcliapp/repository/login"
	repositorylogout "aiagentcliapp/repository/logout"
	repositoryme "aiagentcliapp/repository/me"
	repositoryposition "aiagentcliapp/repository/position"
	repositorytokenstore "aiagentcliapp/repository/tokenstore"
	servicedepartment "aiagentcliapp/service/department"
	serviceemployee "aiagentcliapp/service/employee"
	servicelogin "aiagentcliapp/service/login"
	inputbuilder "aiagentcliapp/service/login/input/builder"
	servicelogout "aiagentcliapp/service/logout"
	serviceme "aiagentcliapp/service/me"
	serviceposition "aiagentcliapp/service/position"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	// 終了コードは、CLIを呼び出したシェルやAIが結果を機械的に判断するための値
	ExitOK      = 0
	ExitRuntime = 1
	ExitUsage   = 2
)

// io.Writerのような型が無いため独自の関数の型をつくる
type PasswordReader func(fd int) ([]byte, error)

type accessTokenStore interface {
	SaveAccessToken(context.Context, string) error
	GetAccessToken(context.Context) (string, error)
	DeleteAccessToken(context.Context) error
}

var newTokenStoreRepository = func() accessTokenStore {
	return repositorytokenstore.NewRepository()
}

func Run(
	args []string,
	client *repository.Client,
	ctx context.Context,
	stdin *os.File,
	stdout, stderr io.Writer,
	passwordReader PasswordReader,
) int {
	// サブコマンドが無い場合は使い方の誤りとして扱い、使用方法をstderrへ出す。
	if len(args) == 0 {
		printUsage(stderr)
		return ExitUsage
	}

	// 先頭の引数をサブコマンド名として扱い、残りは各コマンドの引数にする。
	cmd := args[0]

	// 入力ミスと実行中の失敗を分けると、呼び出し側が終了コードで原因を判断できる。
	var runtimeErr error
	var validationErr error

	// Controllerはコマンド名を判定し、実際の業務処理はServiceへ委譲する。
	switch cmd {

	case "help", "-h", "--help":
		printHelp(stdout)
		return ExitOK

	case "login":
		// login入力を非表示で行う
		fmt.Fprint(stderr, "Password: ")
		password, err := passwordReader(int(stdin.Fd()))
		// 改行
		fmt.Fprintln(stderr)
		if err != nil {
			runtimeErr = fmt.Errorf("read password: %w", err)
			break
		}

		loginInput, err := inputbuilder.NewBuilder().Build(args[1:], string(password))
		if err != nil {
			validationErr = err
			break
		}

		// テスト容易性のためDIする
		loginRepository := repositorylogin.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		loginService := servicelogin.NewService(loginRepository, tokenStoreRepository)

		runtimeErr = loginService.Login(ctx, loginInput)
		if runtimeErr == nil {
			fmt.Fprintln(stdout, "login succeeded")
		}

	case "me":
		meRepository := repositoryme.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		meService := serviceme.NewService(meRepository, tokenStoreRepository)

		user, err := meService.Me(ctx)
		if err != nil {
			runtimeErr = err
			break
		}
		fmt.Fprintf(stdout, "ID: %d\nName: %s\nEmail: %s\n", user.ID, user.Name, user.Email)

	case "logout":
		tokenStoreRepository := newTokenStoreRepository()
		logoutRepository := repositorylogout.NewRepository(client, tokenStoreRepository)
		logoutService := servicelogout.NewService(logoutRepository)

		runtimeErr = logoutService.Logout(ctx)
		if runtimeErr == nil {
			fmt.Fprintln(stdout, "logout succeeded")
		}

	case "departments":
		searchInput, err := buildDepartmentSearchInput(args[1:])
		if err != nil {
			validationErr = err
			break
		}

		departmentRepository := repositorydepartment.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		departmentService := servicedepartment.NewService(departmentRepository, tokenStoreRepository)

		departments, err := departmentService.Departments(ctx, searchInput)
		if err != nil {
			runtimeErr = err
			break
		}
		for _, department := range departments {
			fmt.Fprintf(stdout, "ID: %d\nCode: %s\nName: %s\n", department.ID, department.Code, department.Name)
		}

	case "employees":
		searchInput, err := buildEmployeeSearchInput(args[1:])
		if err != nil {
			validationErr = err
			break
		}

		employeeRepository := repositoryemployee.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		employeeService := serviceemployee.NewService(employeeRepository, tokenStoreRepository)

		employees, err := employeeService.Employees(ctx, searchInput)
		if err != nil {
			runtimeErr = err
			break
		}
		for _, employee := range employees {
			fmt.Fprintf(
				stdout,
				"ID: %d\nEmployeeNumber: %s\nName: %s %s\nNameKana: %s %s\nEmail: %s\nEmploymentStatus: %s\nDepartment: %d %s %s\nPosition: %d %s %s\n",
				employee.ID,
				employee.EmployeeNumber,
				employee.FamilyName,
				employee.GivenName,
				employee.FamilyNameKana,
				employee.GivenNameKana,
				employee.Email,
				employee.EmploymentStatus,
				employee.Department.ID,
				employee.Department.Code,
				employee.Department.Name,
				employee.Position.ID,
				employee.Position.Code,
				employee.Position.Name,
			)
		}

	case "positions":
		searchInput, err := buildPositionSearchInput(args[1:])
		if err != nil {
			validationErr = err
			break
		}

		positionRepository := repositoryposition.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		positionService := serviceposition.NewService(positionRepository, tokenStoreRepository)

		positions, err := positionService.Positions(ctx, searchInput)
		if err != nil {
			runtimeErr = err
			break
		}
		for _, position := range positions {
			fmt.Fprintf(stdout, "ID: %d\nCode: %s\nName: %s\n", position.ID, position.Code, position.Name)
		}

	default:
		// 想定外commandの場合はエラー出力
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", cmd)
		printUsage(stderr)
		return ExitUsage
	}

	// 入力エラーは使用方法と一緒に表示し、実行時エラーとは別の終了コードにする。
	if validationErr != nil {
		printError(stderr, validationErr)
		printUsage(stderr)
		return ExitUsage
	}

	if runtimeErr != nil {
		printError(stderr, runtimeErr)
		return ExitRuntime
	}

	return ExitOK
}

// printUsageは、入力が足りない・間違っている場合に最低限の使い方を表示する。
func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `Usage:
	login [--email EMAIL] [--token-name NAME]  Login and save access token
	me                                         Show current user
	logout                                     Logout and delete access token
	departments [--order-by id|code|name] [--order-direction asc|desc]
	                                           Search departments
	employees [--keyword KEYWORD] [--department-id ID] [--position-id ID] [--employment-status active|leave|retired]
	                                           Search employees
	positions [--order-by id|code|name] [--order-direction asc|desc]
	                                           Search positions
	help                                       Show command list`)
}

// printHelpは、明示的にhelpが指定されたときのコマンド一覧を表示する。
func printHelp(writer io.Writer) {
	fmt.Fprintln(writer, `Commands:
	login [--email EMAIL] [--token-name NAME]  Login and save access token
	me                                         Show current user
	logout                                     Logout and delete access token
	departments [--order-by id|code|name] [--order-direction asc|desc]
	                                           Search departments
	employees [--keyword KEYWORD] [--department-id ID] [--position-id ID] [--employment-status active|leave|retired]
	                                           Search employees
	positions [--order-by id|code|name] [--order-direction asc|desc]
	                                           Search positions`)
}

// エラー表示を1か所にまとめ、stdoutとstderrの使い分けを崩さないようにする。
func printError(writer io.Writer, err error) {
	fmt.Fprintln(writer, "error:", err)
}

func buildDepartmentSearchInput(args []string) (repositorydepartment.SearchInput, error) {
	flags := flag.NewFlagSet("departments", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	orderBy := flags.String("order-by", "", "order by id, code, or name")
	orderDirection := flags.String("order-direction", "", "order direction asc or desc")
	if err := flags.Parse(args); err != nil {
		return repositorydepartment.SearchInput{}, err
	}
	if flags.NArg() != 0 {
		return repositorydepartment.SearchInput{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if err := validateOneOf("order-by", *orderBy, "id", "code", "name"); err != nil {
		return repositorydepartment.SearchInput{}, err
	}
	if err := validateOneOf("order-direction", *orderDirection, "asc", "desc"); err != nil {
		return repositorydepartment.SearchInput{}, err
	}

	return repositorydepartment.SearchInput{
		OrderBy:        *orderBy,
		OrderDirection: *orderDirection,
	}, nil
}

func buildPositionSearchInput(args []string) (repositoryposition.SearchInput, error) {
	flags := flag.NewFlagSet("positions", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	orderBy := flags.String("order-by", "", "order by id, code, or name")
	orderDirection := flags.String("order-direction", "", "order direction asc or desc")
	if err := flags.Parse(args); err != nil {
		return repositoryposition.SearchInput{}, err
	}
	if flags.NArg() != 0 {
		return repositoryposition.SearchInput{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if err := validateOneOf("order-by", *orderBy, "id", "code", "name"); err != nil {
		return repositoryposition.SearchInput{}, err
	}
	if err := validateOneOf("order-direction", *orderDirection, "asc", "desc"); err != nil {
		return repositoryposition.SearchInput{}, err
	}

	return repositoryposition.SearchInput{
		OrderBy:        *orderBy,
		OrderDirection: *orderDirection,
	}, nil
}

func buildEmployeeSearchInput(args []string) (repositoryemployee.SearchInput, error) {
	flags := flag.NewFlagSet("employees", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	keyword := flags.String("keyword", "", "search keyword")
	departmentID := flags.String("department-id", "", "department id")
	positionID := flags.String("position-id", "", "position id")
	employmentStatus := flags.String("employment-status", "", "employment status active, leave, or retired")
	if err := flags.Parse(args); err != nil {
		return repositoryemployee.SearchInput{}, err
	}
	if flags.NArg() != 0 {
		return repositoryemployee.SearchInput{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if err := validatePositiveInteger("department-id", *departmentID); err != nil {
		return repositoryemployee.SearchInput{}, err
	}
	if err := validatePositiveInteger("position-id", *positionID); err != nil {
		return repositoryemployee.SearchInput{}, err
	}
	if err := validateOneOf("employment-status", *employmentStatus, "active", "leave", "retired"); err != nil {
		return repositoryemployee.SearchInput{}, err
	}

	return repositoryemployee.SearchInput{
		Keyword:          *keyword,
		DepartmentID:     *departmentID,
		PositionID:       *positionID,
		EmploymentStatus: *employmentStatus,
	}, nil
}

func validateOneOf(name, value string, allowed ...string) error {
	if value == "" {
		return nil
	}
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %v", name, allowed)
}

func validatePositiveInteger(name, value string) error {
	if value == "" {
		return nil
	}
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return fmt.Errorf("%s must be a positive integer", name)
	}
	return nil
}
