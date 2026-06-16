package main

import (
        "context"
        "database/sql"
        "errors"
        "example.com/internal/database"
        "fmt"
        "github.com/google/uuid"
        "time"
)

func handlerLogin(s *state, cmd command) error {
        if len(cmd.Args) != 1 {
                return fmt.Errorf("usage: %s <name>", cmd.Name)
        }
        name := cmd.Args[0]

    _, err := s.db.GetUser(context.Background(), name)
    if errors.Is(err, sql.ErrNoRows) {
        return fmt.Errorf("user: %v does not exist", name)
    }
    if err != nil {
        return fmt.Errorf("couldn't check if user exists %w", err)
    }

        err = s.cfg.SetUser(name)
        if err != nil {
                return fmt.Errorf("couldn't set current user: %w", err)
        }

        fmt.Println("User switched successfully!")
        return nil
}

func handlerRegister(s *state, cmd command) error {
        // 1 check if a name was pased to the command len(cmd.Args)
        if len(cmd.Args) != 1 {
                return fmt.Errorf("usage: %s <name>", cmd.Name)
        }

        // 2 name := cmd.Args[0]
        name := cmd.Args[0]

        // check if the user already exists
        _, err := s.db.GetUser(context.Background(), name)
        if err == nil {
                // user exists
                return fmt.Errorf("user: %s already exists.", name)
        }

        if !errors.Is(err, sql.ErrNoRows) {
                return fmt.Errorf("couldn't check if user exists: %w", err)
        }

        // 3 call s.db.CreateUser(context.Background(), database.CreateUserParams{...})
        user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
                ID:        uuid.New(),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Name:      name,
        })

        // 4 if error, return it
        if err != nil {
                return fmt.Errorf("couldn't register new users: %w", err)
        }

        // 5 set current user in config
        err = s.cfg.SetUser(name)
        if err != nil {
                return fmt.Errorf("couldn't set current user: %w", err)
        }

        // 6 print success/debug info
        fmt.Printf("User successfully created: %v\n", user)

        // 7 return nil
        return nil
}

func handlerReset(s *state, cmd command) error {
    // do some resetting baby
    if len(cmd.Args) > 0 {
        return fmt.Errorf("usage: %s", cmd.Name)
    }

    err := s.db.Reset(context.Background())
    if err != nil {
        return fmt.Errorf("unable to reset database: %w", err)
    }

    fmt.Println("database reset")
    return nil
}

func handlerUsers(s *state, cmd command) error {
    // look for some users
    if len(cmd.Args) > 0 {
        return fmt.Errorf("usage: %s", cmd.Name)
    }

    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return fmt.Errorf("unable to get a list of users: %w", err)
    }

    for _, user := range users {
        if user.Name == s.cfg.CurrentUserName {
            fmt.Printf("* %v (current)\n", user.Name)
        } else {
            fmt.Printf("* %v\n", user.Name)
        }
    }

    return nil
}

