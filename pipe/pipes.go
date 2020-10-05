package pipe

//PipeFromExecutable creates a plumber.Pipe that uses given executable
//func PipeFromExecutable(path string, needsState bool) plumber.Pipe {
//	return func(s plumber.State, input interface{}) (interface{}, error) {
//		all, err := s.All()
//		if err != nil {
//			return nil, err
//		}
//		bs, err := json.Marshal(all)
//		if err != nil {
//			return nil, err
//		}

//		args := []string{}
//		if needsState {
//			args = append(args, fmt.Sprintf("\"%s\"", string(bs)))
//		}
//		args = append(args, fmt.Sprintf("\"%v\"", input))

//		c := exec.Command(path, args...)

//		output, err := c.Output()
//		if err != nil {
//			return nil, err
//		}
//		//TODO: should update state cause based on the contract updated fileds are in output
//		return string(output), nil
//	}
//}
