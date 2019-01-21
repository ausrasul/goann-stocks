#!/usr/bin/perl
#
# Dra ut en 30 dagars time series på close där close' = close'/close[idag]-1
# Dra sedan ut outputs som indikerar +-1% inom 5 dagar, +-2% inom 5 dagar, +-3% inom 5 dagar, +-4% inom 5 dagar
#

if ($ARGV[0] eq '-p') {
	$userPassword = "-p$ARGV[1]";
	splice(@ARGV,0,2);
}
if ($ARGV[0] eq '-f') {
	$format = $ARGV[1];
	splice(@ARGV,0,2);
}
$commodity = shift @ARGV;
$sdt = shift @ARGV;
$edt = shift @ARGV;
@averages = @ARGV;

$id = `mysql -uroot $userPassword borsa -N -B -e 'select id from commodity where shortname = "$commodity"'`;
chomp($id);
die("No such commodity: $commodity\n") unless $id =~ /\d+/;

$maxpoints = 40;
for (@averages) { if ($_ > $maxpoints) { $maxpoints = $_; } }

print "{ \"values\": [\n" if $format eq 'json';
open(PRICE, "mysql -uroot $userPassword borsa -N -B -e 'select dt,close,high,low from commodity_price where id=$id and dt>=\"$sdt\" and dt<=\"$edt\" order by dt asc'|");
while ($l = <PRICE>) {
	chomp($l);
	($dt,$close,$high,$low) = split /\t/, $l;
	push @dt, $dt;
	push @close, $close;
	push @high, $high;
	push @low, $low;
	if ($#close+1 > $maxpoints) {
		shift @dt;
		shift @close;
		shift @high;
		shift @low;

		my @tmp;
		#@tmp = @dt[-10..-5]; # 6 dagar -- 2007-01-10 2007-01-11 2007-01-12 2007-01-15 2007-01-16 2007-01-17
		#@tmp = @dt[-4..-1]; # 4 dagar -- 2007-01-18 2007-01-19 2007-01-22 2007-01-23
		# Ovan ger: Future = -1 -> -5
		#           Today  = -6
		#           Past   = -7 -> -36

		$dt_today = $dt[-6];
		$close_today = $close[-6];
		$high_today = $high[-6];
		$low_today = $low[-6];

		# Next 5 days high
		@tmp = @high[-5..-1];
		$n5one = 0;
		for $tmp (@tmp) {
			if ($tmp >= $close_today * 1.01) {
				$n5one = 1;
				last;
			}
		}
		$n5two = 0;
		for $tmp (@tmp) {
			if ($tmp >= $close_today * 1.02) {
				$n5two = 1;
				last;
			}
		}
		$n5three = 0;
		for $tmp (@tmp) {
			if ($tmp >= $close_today * 1.03) {
				$n5three = 1;
				last;
			}
		}

		# Next 5 days low
		@tmp = @low[-5..-1];
		$n5lone = 0;
		for $tmp (@tmp) {
			if ($tmp <= $close_today * 0.99) {
				$n5lone = 1;
				last;
			}
		}
		$n5ltwo = 0;
		for $tmp (@tmp) {
			if ($tmp <= $close_today * 0.98) {
				$n5ltwo = 1;
				last;
			}
		}
		$n5lthree = 0;
		for $tmp (@tmp) {
			if ($tmp <= $close_today * 0.97) {
				$n5lthree = 1;
				last;
			}
		}

		# Calculate the timeseries into the range 0 to 1
		# Let 0.9 => 0 and 1.1 be 1, any value outside will be scaled to 0 or 1
		my @closepct;
		for $tmp (@close[-36..-7]) {
            if ($close_today == 0){
                $x = 0;
            } else {
			    $x = sprintf("%.2f", ($tmp / $close_today));
            }
			$x -= 0.9;
			if ($x < 0) {
				$x = 0;
			} elsif ($x > 0.2) {
				$x = 1;
			} else {
				$x = sprintf("%.2f", ($x / 0.2));
			}
			push @closepct, $x;
		}

		# Print
		print join("\t", @closepct) . "\t$n5one\t$n5two\t$n5three\t$n5lone\t$n5ltwo\t$n5lthree\n" if $format eq 'ft';
		print "$dt_today\t$close_today\t$high_today\t$low_today\t" . join("\t", @closepct) . "\t$n5one\t$n5two\t$n5three\t$n5lone\t$n5ltwo\t$n5lthree\n" if $format eq 'fte';
		if ($firstrow++) { print ",\n" if $format eq 'json'; }
		print "{ \"date\": \"$dt_today\", \"close\": $close_today, \"high\": $high_today, \"low\": $low_today,\n\"inputs\": [" . join(",", @closepct) . "],\n\"outputs\": [$n5one,$n5two,$n5three,$n5lone,$n5ltwo,$n5lthree]}" if $format eq 'json';
	}
}
close(PRICE);
print "\n]}\n" if $format eq 'json';
