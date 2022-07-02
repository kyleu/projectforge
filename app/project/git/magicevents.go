package git

func (s *Service) onAhead(args *magicArgs, add func(msg string, args ...any)) error {
	if args.Dirty == 0 {
		args.Result.Status = "push"
	} else {
		args.Result.Status = "commit, push"
		if err := s.magicCommit(args, add); err != nil {
			return err
		}
	}
	if err := s.magicPush(args, args.Ahead+1, add); err != nil {
		return err
	}
	return nil
}

func (s *Service) onBehind(args *magicArgs, add func(msg string, args ...any)) error {
	if args.Dirty > 0 {
		if err := s.magicStash(args, add); err != nil {
			return err
		}
	}
	args.Result.Status = "pull"
	if err := s.magicPull(args, add); err != nil {
		return err
	}
	if args.Dirty > 0 {
		args.Result.Status += ", commit"
		if err := s.magicStashPop(args, add); err != nil {
			return err
		}
		if err := s.magicCommit(args, add); err != nil {
			return err
		}
		if err := s.magicPush(args, 1, add); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) onClean(args *magicArgs, add func(msg string, args ...any)) error {
	if args.Dirty > 0 {
		args.Result.Status = "commit"
		if err := s.magicCommit(args, add); err != nil {
			return err
		}
		if err := s.magicPush(args, 1, add); err != nil {
			return err
		}
	}
	return nil
}
