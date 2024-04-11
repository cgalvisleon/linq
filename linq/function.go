package linq

// Count function to use in linq
func (l *Linq) Count(col *Column) *Linq {
	sel := l.GetColumn(col)
	sel.TypeFunction = TpCount

	return l
}

// Sum function to use in linq
func (l *Linq) Sum(col *Column) *Linq {
	sel := l.GetColumn(col)
	sel.TypeFunction = TpSum

	return l
}

// Avg function to use in linq
func (l *Linq) Avg(col *Column) *Linq {
	sel := l.GetColumn(col)
	sel.TypeFunction = TpAvg

	return l
}

// Max function to use in linq
func (l *Linq) Max(col *Column) *Linq {
	sel := l.GetColumn(col)
	sel.TypeFunction = TpMax

	return l
}

// Min function to use in linq
func (l *Linq) Min(col *Column) *Linq {
	sel := l.GetColumn(col)
	sel.TypeFunction = TpMin

	return l
}
