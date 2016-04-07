with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

procedure exercise8 is

    Count_Failed    : exception;    -- Exception to be raised when counting fails
    Gen             : Generator;    -- Random number generator

    protected type Transaction_Manager (N : Positive) is
        entry Finished;
		entry Wait_Until_Aborted;

        procedure Signal_Abort;
    private
        Finished_Gate_Open  : Boolean := False;
        Aborted             : Boolean := False;
    end Transaction_Manager;
    protected body Transaction_Manager is
        entry Finished when Finished_Gate_Open or Finished'Count = N is
        begin
            ------------------------------------------
            -- PART 3: Complete the exit protocol here
            ------------------------------------------
			if Finished'Count = N-1 then -- First one in let all in
				Finished_Gate_Open := True;
			end if;

			if Finished'Count = 0 then -- Last one close the door
				Finished_Gate_Open := False;
			end if;
        end Finished;

		------------------------------------------
		-- Exercise 8: PART 2: Create the entry
		------------------------------------------
		entry Wait_Until_Aborted when Aborted is
			count_worker : Integer := 0;
		begin
			count_worker = count_worker +1;
			if count_worker = 3 then -- Reset when last one is finished
				Aborted := False;
				count_worker := 0;
			end if;
		end Wait_Until_Aborted;

        procedure Signal_Abort is
        begin
            Aborted := True;
        end Signal_Abort;

    end Transaction_Manager;


    function Unreliable_Slow_Add (x : Integer) return Integer is
    Error_Rate : Constant := 0.15;  -- (between 0 and 1)
		return_value : Integer := 0;
    begin
        -------------------------------------------
        -- PART 1: Create the transaction work here
        -------------------------------------------
		-- compare Random(Gen) with Error_Rate and do:
		if Random(Gen) > Error_Rate then
			---- The intended behaviour:
			delay Duration(Random(Gen)*4.0); -- Work takes 4-ish seconds
			return_value := x + 10;
			return return_value;
		else
			---- The faulty behaviour:
			delay Duration(Random(Gen)*0.5); -- Work takes up to half a second
			Put_Line("-- ** Exception was raised **");
			raise Count_Failed; -- Error, raise exception
		end if;
    end Unreliable_Slow_Add;

    task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
    task body Transaction_Worker is
        Num         : Integer   := Initial;
        Prev        : Integer   := Num;
        Round_Num   : Integer   := 0;
    begin
        Put_Line ("Worker" & Integer'Image(Initial) & " started");

        loop
            Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
            Round_Num := Round_Num + 1;

			------------------------------------------
			-- Exercise 8: PART 1: Select-Then-Abort
			------------------------------------------
			select
				Manager.Wait_Until_Aborted;  -- eg. X.Entry_Call;
				-- code that is run when the triggering_alternative has triggered
				--   (forward ER code goes here)
				Num := Num + 5;
			then abort
				begin
					Num := Unreliable_Slow_Add (Num); -- Add Num +10
				exception -- Start of exception handlers
					when Count_Failed =>
					    Put_Line("-- Exceptiopn Worker " & Integer'Image(Initial));
						Manager.Signal_Abort;
				end;
			end select;

			Manager.Finished;
			Put_Line ("-- Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num));

            Prev := Num;
            delay 0.5;

        end loop;

    end Transaction_Worker;

    Manager : aliased Transaction_Manager (3);

    Worker_1 : Transaction_Worker (0, Manager'Access);
    Worker_2 : Transaction_Worker (1, Manager'Access);
    Worker_3 : Transaction_Worker (2, Manager'Access);

begin
    Reset(Gen); -- Seed the random number generator
end exercise8;
