HARDWARE DEBUG/NCD DEBUG
    INPUT
        1) Temp directory [/tmp/Raju/ will get deleted with time. So 
            create this folder
        2) Input file should be with *.txt extension and the name better
            should be like appserv93.txt 
        3) If using different systems, rename each file with it's name 
            and copy it there
        4) ncdutils -a and copy the output into this files
    RUN
        ./hwdiag -path=/Users/rajuv/Techsupport/6477
    BUILD
        call ./build script. This will compile all the required files
    TODO
        1) Statistics [ncdutil -s]
PKT DECODE
    INPUT: 
        1) Copy the hexa captured pkt output to InputFiles/pktinput.txt
    RUN
        1) go run pktDecode.go
COUNTERS DECODE
    ncdutils -s # Dump statistics
    Drop Events : Pkts dropped due to Insufficient Buffer Overrun. 
        Flow Control on ingress port. no enough buffer to store the 
        incoming packets
    Oversize Packets: Pkts dropped due to receving Oversized Packets
    MAC Rx Err Pkt Rcvd: ??
    Bad FC Received: ??
    Deferred Pkt Sent: ??
    Bad CRC:
    Babber Packets
    BRDC Packets Setn
    Excessive Collision
    MAC Transmit Error
    Bad Octets Received:
    Undersize Packets:
    Multiple Pkt Sent:
    Collisions:
    Late Collisions
STATISTICS
    Statistics Info:
       =======================
       nm_hello: 1, nm_query_portcfg: 1, nm_heartbeat: 50,nm_query_switchcfg: 1, nm_query_boardinfo: 1, nm_config_switchif: 4, nm_disconnect: 0
       nm_vtep_create: 0, nm_vtep_delete: 0, nm_create_local_if: 5, nm_delete_local_if: 0, nm_query_local_if: 2, nm_create_remote_if: 0
       nm_delete_remote_if: 0, nm_create_remote_vtep: 0, nm_delete_remote_vtep: 0, nm_enable_vlan: 2, nm_disable_vlan: 3, nm_query_vlan: 2
       nm_enable_vxlan: 0, nm_disable_vxlan: 0, nm_add_entry_vxlan: 0, nm_remove_entry_vxlan: 0, nm_add_qostemplate: 3, nm_remove_qostemplate: 0
       nm_update_qostemplate: 0, nm_query_network_stats: 0 conduit_timeout_cnt: 47805

       nm_hello : This must be 1. i.e Bosun sends Hello message only one 
            time during bootup and Embedded will read the file to get FC VLAN 
            etc
        nm_disconnect : this must be 0. Disconnect mesage from Bosun 
            is received. During this hb_state, will be set to 0 and 
            also last_uptime will be set 0
        nm_heartbeat : This is incremented when nCD received heart beat
            message. This should keep getting incremented
    curr_uptime: 990, last_uptime: 986, hb_state: 1, bcast_mcast_ntwk_pif: 0x8
    hb_state : This must be 1. Ie Healthy connection
    last_uptime and curr_uptime. The different b/w the times < 90
        NCD Logs :  - Check for this Error
        SKIPPER_LOG_ERR("Heartbeat missed at: %ld [last recorded hb: %ld]\n",
                        curr_uptime, skipper_cfg_info.hb_uptime);
DESIGN
    TextProcessingTool.pdf
    FDB Database
        key:Mac,vlan 
        Value :  Map[Gport] Map of Gports
    Port Database
        Key:server,port
        value : map[MacVlan]
    PIN Database
        key:server,port
        value: server,port

DEBUG LIVE SWITCH
    ./rajuDiag.sh appserv93 appserv94
    # don't run it as rajuDiag.sh app..use ./
    1) This will run diah.sh on each of the remote servers. This check 
        for status of Services and collect logs (bosun,convoy and armada) 
        and sort those files by time into single file [ rajumerged.txt ]
    2) This create rajumerged.txt file in each of the server /home/diamanti
    3) that file is copied to this folder as appserv93Raju.txt..etc
    4) Embedded Commnand output is copied to appserv93RajuEmbed.txt
