package ng

import (
	"context"
	"github.com/bombsimon/logrusr/v4"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestSolver_Solve(t *testing.T) {
	logger := logrusr.New(logrus.StandardLogger())
	solver, err := NewHungerSolver(SolverOptions{
		Logger: logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	compass, err := ParseCompass("0+3,3-3,0+2/mi,mo,io")
	if err != nil {
		t.Fatal(err)
	}

	solution, err := solver.Solve(context.Background(), compass)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(compass.String())
	t.Log(solution.String())
}
