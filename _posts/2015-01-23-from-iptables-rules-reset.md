---
layout: post
title:  "How to reset all iptables rules configured"
date:   2015-01-23 10:15:44
categories: TA
---

If you ask somebody how to reset all iptables rules configured, he or she may answer you in 3 seconds. The answer may be:

    iptables -F

With the help document of iptables, it says:

      --flush   -F [chain]		Delete all rules in  chain or all chains

For me, when I read the sentences above, I think I got the point.

A few days ago, I wanted to forward HTTP service from localhost to a remote server, I used the following command:

    iptables -t nat -I PREROUTING -p tcp --dport 80 -j DNAT --to <remote ip>

Then I executed the command `iptables -L`, but I did not get the new rule.

    Chain INPUT (policy ACCEPT)
    target     prot opt source               destination

    Chain FORWARD (policy ACCEPT)
    target     prot opt source               destination

    Chain OUTPUT (policy ACCEPT)
    target     prot opt source               destination

At this moment, I was told the remove server IP address had been changed. So I wanted to flush the current rule and configured a new one.

But after executing `iptables -F`, the forwarding still worked.

Why? After searching some meterials and reading the manual carefully, I find the point finally.

`iptables` has a parameter `-t`, it will specify a table, defaultly the table is `filter`. But in the command above I configured a rule under table `nat`. So `iptables -F` will not flush that rule, as it does not belong the `filter`.

Now I execute `iptables -t nat -L`, I get it.

    Chain PREROUTING (policy ACCEPT)
    target     prot opt source               destination
    DNAT       tcp  --  anywhere             anywhere            tcp dpt:hbci to:<dest ip>
    DNAT       tcp  --  anywhere             anywhere            tcp dpt:http to:<dest ip>

    Chain POSTROUTING (policy ACCEPT)
    target     prot opt source               destination
    MASQUERADE  all  --  anywhere             anywhere

    Chain OUTPUT (policy ACCEPT)
    target     prot opt source               destination

Then `iptables -t nat -F` remove the rule I just configured.

So the solution to remove all rules is:

    iptables -F
    iptables -X
    iptables -t nat -F
    iptables -t nat -X
    iptables -t mangle -F
    iptables -t mangle -X
    iptables -P INPUT ACCEPT
    iptables -P FORWARD ACCEPT
    iptables -P OUTPUT ACCEPT
