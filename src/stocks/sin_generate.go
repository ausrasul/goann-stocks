package main
import "fmt"
//import "log"
import "math/rand"

func main(){
    fmt.Println("1000 6 1")
    var amp float64 = 0
    s1 := rand.NewSource(42)
    r1 := rand.New(s1)
	for i:=0;i<1000;i++{
		for j:=0;j<6;j++{
			amp = r1.Float64()// math.Sin((float64(i+j)/10))
			fmt.Print(amp, " ")
		}
		fmt.Println()
		namp := r1.Float64()//math.Sin((float64(i+6)/10))
        if namp > amp {
		    fmt.Println(1)
        } else {
            fmt.Println(0)
        }
	}
}
