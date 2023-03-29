package orm

// 更新的时候

/*
xxxUpdateAll 更新req struct中对应的所有字段（包括零语义）
	eg:StateMachineUpdateAll  req struct.SM所有字段需要更新

yyyUpdateS[other word] 除了函数备注的字段外，只更新req struct中非零语义

zzzUpdate[other word] 只更新req struct中非零语义

既有零语义又有非零语义，增加额外的标识来确定
*/
