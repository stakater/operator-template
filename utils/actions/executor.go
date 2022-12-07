package actions

type ActionQueue struct {
	Actions []IAction
}

func (ex *ActionQueue) Enqueue(a IAction) {
	if ex.Actions == nil {
		ex.Actions = []IAction{}
	}

	ex.Actions = append(ex.Actions, a)
}

func (ex *ActionQueue) Run() (IAction, error) {
	for _, action := range ex.Actions {
		err := action.Run()
		if err != nil {
			return action, err
		}
	}

	return nil, nil
}
