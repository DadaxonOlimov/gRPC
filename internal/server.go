package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"github.com/DadaxonOlimov/model"
	"github.com/DadaxonOlimov/proto"
)

type grpcServer struct {
	proto.UnimplementedTodoServiceServer
	db *sql.DB
}

func (g grpcServer) Create(ctx context.Context, request *proto.CreateTodoRequest) (*proto.CreateTodoResponse, error) {
	var newTodo proto.Todo
	query := g.db.QueryRow("insert into todos (user_id, description, title) values ($1, $2, $3) returning *",
		request.UserId, request.Description, request.Title)
	err := query.Scan(&newTodo.Id, &newTodo.UserId, &newTodo.Description, &newTodo.Title)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &proto.CreateTodoResponse{
		Id:          newTodo.Id,
		UserId:      newTodo.UserId,
		Description: newTodo.Description,
		Title:       newTodo.Title,
	}, nil
}

func (g grpcServer) GetById(ctx context.Context, request *proto.GetByIdRequest) (*proto.GetTodoResponse, error) {
	//TODO implement me
	var todo model.Todos
	rows, err := g.db.Query("select * from todos where id = $1", request.Id)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		return nil, err
	}
	rows.Next()
	if err := rows.Scan(&todo.Id, &todo.UserId, &todo.Description, &todo.Title); err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		return nil, err
	}
	fmt.Println(todo.Id)
	return &proto.GetTodoResponse{
		Id:          fmt.Sprintf("%d", todo.Id),
		UserId:      todo.UserId,
		Description: todo.Description,
	}, nil
}

func (g grpcServer) GetAll(ctx context.Context, request *proto.GetAllTodosRequest) (*proto.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	rows, err := g.db.Query(`SELECT * FROM todos`)
	if err != nil {
		return nil, err
	}
	var allTodos []*proto.Todo

	for rows.Next() {
		var newTodo proto.Todo
		err := rows.Scan(&newTodo.Id, &newTodo.UserId, &newTodo.Description, &newTodo.Title)
		if err != nil {
			fmt.Println(fmt.Errorf("An error occured: %s ", err.Error()))
			return nil, err
		}

		allTodos = append(allTodos, &newTodo)
	}

	var todosLen = len(allTodos)

	if todosLen == 0 {
		return &proto.GetAllResponse{
			Todos: []*proto.Todo{},
		}, nil
	}

	return &proto.GetAllResponse{
		Todos: allTodos,
	}, nil
}

func (g grpcServer) Delete(ctx context.Context, requset *proto.DeleteTodoRequset) (*proto.DeleteTodoResponse, error) {
	fmt.Println(requset.Id, "deleted ID")
	_, err := g.db.Exec("delete from todos where id = $1", requset.Id)
	if err != nil {
		panic(err)
	}
	return &proto.DeleteTodoResponse{Id: requset.Id}, nil
}

func (g grpcServer) Update(ctx context.Context, request *proto.UpdateTodoRequest) (*proto.UpdateTodoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	_, err := g.db.Exec("update todos set user_id = $1, title = $2,description = $3 where id = $4", request.UserId, request.Description, request.Title, request.Id)
	if err != nil {
		panic(err)
	}
	return &proto.UpdateTodoResponse{
		Id: request.Id,
	}, nil
}

func NewGRPCServer(db *sql.DB) grpcServer {
	return grpcServer{
		db: db,
	}
}