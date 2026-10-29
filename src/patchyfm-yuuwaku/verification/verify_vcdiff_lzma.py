from pathlib import Path
import lzma, zlib, hashlib, sys, struct

srcp=Path(sys.argv[1]); patchp=Path(sys.argv[2]); outp=Path(sys.argv[3]);
source=srcp.read_bytes(); p=patchp.read_bytes(); pos=0

def intat(buf, index):
    v=0
    while True:
        c=buf[index];index+=1
        v=(v<<7)|(c&127)
        if not c&128:return v,index

def readvar():
    global pos
    value,pos=intat(p,pos)
    return value
assert p[:4]==b'\xd6\xc3\xc4\x00'
pos=4
hd=p[pos];pos+=1
if hd&1:
    compressor=p[pos];pos+=1
    assert compressor==2,compressor
else: compressor=None
assert not hd&2
if hd&4:
    n=readvar(); app=p[pos:pos+n];pos+=n
secdec=[lzma.LZMADecompressor(format=lzma.FORMAT_XZ) for _ in range(3)]
table=[[(2,0,0)]]
ct=[None]*256
ct[0]=[(0,0,0)]
for op in range(1,19):ct[op]=[(1,0 if op==1 else op-1,0)]
for op in range(19,163):
    q=op-19;ct[op]=[(2, 0 if q%16==0 else q%16+3,q//16)]
for op in range(163,235):
    q=op-163; mode=q//12;add_size=q%12//3+1; copy_size=q%3+4
    ct[op]=[(1,add_size,0),(2,copy_size,mode)]
for op in range(235,247):
    q=op-235;ct[op]=[(1,q%4+1,0),(2,4,6+q//4)]
for op in range(247,256):ct[op]=[(2,4,op-247),(1,1,0)]
assert all(ct)

with outp.open('wb') as outf:
    produced=0
    winidx=0
    while pos<len(p):
        st=pos; w=p[pos];pos+=1
        assert w&0b11==1,('unsupported window',w)
        sl=readvar();so=readvar()
        dl=readvar();tl=readvar();ind=p[pos];pos+=1
        lengths=[readvar() for _ in range(3)]
        chk=struct.unpack('>I',p[pos:pos+4])[0] if w&4 else None
        if w&4:pos+=4
        segments=[]
        for t,n in enumerate(lengths):
            dat=p[pos:pos+n];pos+=n
            if ind&(1<<t):
                expected,off=intat(dat,0)
                b=secdec[t].decompress(dat[off:])
                assert len(b)==expected,(winidx,t,expected,len(b))
            else:b=dat
            segments.append(b)
        dat,inst,addr=segments
        print('window',winidx,'source',sl,so,'target',tl,'raw sec',*[len(q) for q in segments],'comp',lengths, 'adler',hex(chk) if chk else None,flush=True)
        result=bytearray();ni=na=nd=0
        near=[0]*4; same=[0]*768; last=0
        source_region=source[so:so+sl]
        assert len(source_region)==sl
        count=0
        while ni<len(inst):
            opcode=inst[ni];ni+=1
            for typ, sz, mode in ct[opcode]:
                if sz==0:
                    sz,ni=intat(inst,ni)
                if typ==0: # RUN
                    b=dat[nd];nd+=1;result.extend(bytes([b])*sz)
                elif typ==1: # ADD
                    result.extend(dat[nd:nd+sz]);nd+=sz
                else:
                    here=sl+len(result)
                    if mode==0:
                        a,na=intat(addr,na)
                    elif mode==1:
                        t,na=intat(addr,na);a=here-t
                    elif mode<=5:
                        t,na=intat(addr,na);a=near[mode-2]+t
                    else:
                        a=same[(mode-6)*256+addr[na]];na+=1
                    near[last]=a;last=(last+1)%4;same[a%768]=a
                    if a<sl:
                        n=min(sz,sl-a)
                        result+=source_region[a:a+n]
                        sz-=n; a+=n
                    if sz:
                        a-=sl
                        for k in range(sz):
                            result.append(result[a+k])
                count+=1
        assert len(result)==tl,(winidx,len(result),tl)
        assert na==len(addr),(winidx,na,len(addr))
        assert nd==len(dat),(winidx,nd,len(dat))
        if chk is not None:
            got=zlib.adler32(result)&0xffffffff
            assert got==chk,('Adler mismatch',winidx,hex(got),hex(chk))
        outf.write(result);produced+=len(result);winidx+=1
print('all bytes',produced,'md5',hashlib.md5(outp.read_bytes()).hexdigest(),'sha256',hashlib.sha256(outp.read_bytes()).hexdigest())
