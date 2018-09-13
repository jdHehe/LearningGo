package NeuralNetwork

// 神经网络的BP算法的golang实现

// 简单的三层神经网络：输入层、隐藏层、输出层
// 变量定义： 输入层InputLayers := o
//            隐藏层HiddenLayer := p
//            输出层OutputLayer := q
//            输入层到隐藏层的网络参数：v=(o, p) 括号内为矩阵的shape属性
//            隐藏层的阈值 gamma=(p,)
//            隐藏层到输出层的参数: w=(p,q)
//            输出层的阈值 theta=(q,)

type FeedForward struct {
	InputLayer, HiddenLayer, OutputLayer int
	V			[][]float64
	Gamma 		[]float64
	W 			[][]float64
	Theta 		[]float64
}

func (ff *FeedForward)Init(input int, hidden int, output int){
	ff.InputLayer = input
	ff.OutputLayer = output
	ff.HiddenLayer = hidden
	ff.V = makeMatrix(input, hidden)
	ff.Gamma = make([]float64, hidden)
	ff.W = makeMatrix(hidden, output)
	ff.Theta = make([]float64, output)
}

func (ff *FeedForward)Train_BP(learning_rate float64){
	// 下面两个分别是  隐藏层的输出 和 输出层的输出
	HiddenUnitValue := make([]float64, ff.HiddenLayer) //h_i
	OutputUnitValue := make([]float64, ff.OutputLayer) //y^_i
	g := make([]float64, ff.OutputLayer) //输出层到隐藏层的梯度
	e := make([]float64, ff.HiddenLayer) //隐层层到输入层的梯度

	// 前向和后向反馈
	forward := func(f *FeedForward, X[]float64, y[]float64) {
		for i:=0; i<len(HiddenUnitValue); i++{
			var tmp float64 = 0
			for j:=0; j<f.InputLayer; j++{
				tmp += ( f.V[j][i] * X[j] )
			}
			tmp -= ff.Gamma[i]
			HiddenUnitValue[i] = SigmodFunction(tmp)
		}

		for i:=0; i<len(OutputUnitValue); i++{
			var tmp float64 = 0
			for j:=0; j<f.HiddenLayer; j++{
				tmp += ( f.W[j][i] * HiddenUnitValue[j] )
			}
			tmp -= ff.Theta[i]
			OutputUnitValue[i] = SigmodFunction(tmp)
		}
	}

	backward := func(f *FeedForward, X[]float64, y[]float64) {
		//  计算两层的梯度值g和e
		for i:=0; i<f.OutputLayer; i++{
			g[i] = OutputUnitValue[i] * (1-OutputUnitValue[i]) * (y[i] - OutputUnitValue[i])
		}

		for i:=0; i<f.HiddenLayer; i++{
			var tmp float64 = 0
			for j:=0; j<f.OutputLayer; j++{
				tmp  += (g[j]*f.W[i][j])
			}
			e[i] = HiddenUnitValue[i] * (1-HiddenUnitValue[i]) * tmp
		}

		//	根据梯度更新  神经网络的参数 以及阈值
		for i:=0; i<f.HiddenLayer;i++  {
			for j:=0; j<f.OutputLayer;j++  {
				f.W[i][j] += learning_rate*g[j]*HiddenUnitValue[i]
			}
			f.Gamma[i] += -(learning_rate * e[i])
		}
		for i:=0; i<f.InputLayer; i++{
			for j:=0; j<f.HiddenLayer; j++{
				f.V[i][j] += learning_rate*e[j]*X[i]
			}
		}
		for i:=0; i<f.OutputLayer; i++{
			f.Theta[i] += -(learning_rate*g[i])
		}
	}


}