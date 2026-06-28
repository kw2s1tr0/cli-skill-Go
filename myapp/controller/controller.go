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
	departmentinputbuilder "aiagentcliapp/service/department/input/builder"
	serviceemployee "aiagentcliapp/service/employee"
	employeeinputbuilder "aiagentcliapp/service/employee/input/builder"
	servicelogin "aiagentcliapp/service/login"
	logininputbuilder "aiagentcliapp/service/login/input/builder"
	servicelogout "aiagentcliapp/service/logout"
	serviceme "aiagentcliapp/service/me"
	serviceposition "aiagentcliapp/service/position"
	positioninputbuilder "aiagentcliapp/service/position/input/builder"
	"context"
	"fmt"
	"io"
	"os"
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

		loginInput, err := logininputbuilder.NewBuilder().Build(args[1:], string(password))
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
		searchInput, err := departmentinputbuilder.NewBuilder().Build(args[1:])
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
		searchInput, err := employeeinputbuilder.NewBuilder().Build(args[1:])
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
		searchInput, err := positioninputbuilder.NewBuilder().Build(args[1:])
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
